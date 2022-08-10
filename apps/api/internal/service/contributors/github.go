package contributors

import (
	"context"
	"net/http"

	"contrib.rocks/apps/api/internal/tracing"
	"contrib.rocks/libs/goutils/model"
	"github.com/google/go-github/v45/github"
	"golang.org/x/sync/errgroup"
)

func fetchRepositoryContributors(client *github.Client, c context.Context, r *model.Repository) (*model.RepositoryContributors, error) {
	ctx, span := tracing.Tracer().Start(c, "contributors.fetchRepositoryContributors")
	defer span.End()

	repositoryResultOut := make(chan *github.Repository, 1)
	contributorsResultOut := make(chan []*github.Contributor, 1)
	eg, ctx := errgroup.WithContext(ctx)
	eg.Go(makeFetchRepositoryFn(client, ctx, r, repositoryResultOut))
	eg.Go(makeFetchContributorsFn(client, ctx, r, contributorsResultOut))

	if err := eg.Wait(); err != nil {
		return nil, err
	}
	repositoryResult := <-repositoryResultOut
	contributorsResult := <-contributorsResultOut

	contributors := make([]*model.Contributor, 0, len(contributorsResult))
	for _, e := range contributorsResult {
		contributors = append(contributors, &model.Contributor{
			ID:            e.GetID(),
			Login:         e.GetLogin(),
			AvatarURL:     e.GetAvatarURL(),
			HTMLURL:       e.GetHTMLURL(),
			Contributions: e.GetContributions(),
		})
	}
	return &model.RepositoryContributors{
		Repository: &model.Repository{
			Owner:    repositoryResult.Owner.GetLogin(),
			RepoName: repositoryResult.GetName(),
		},
		StargazersCount: repositoryResult.GetStargazersCount(),
		Contributors:    contributors,
	}, nil
}

func makeFetchRepositoryFn(client *github.Client, c context.Context, repo *model.Repository, out chan<- *github.Repository) func() error {
	return func() error {
		data, resp, err := client.Repositories.Get(c, repo.Owner, repo.RepoName)
		if err != nil {
			if resp != nil && resp.StatusCode == http.StatusNotFound {
				return &RepositoryNotFoundError{repo}
			}
			return err
		}
		if out != nil && cap(out) > 0 {
			out <- data
		}
		return nil
	}
}

func makeFetchContributorsFn(client *github.Client, c context.Context, repo *model.Repository, out chan<- []*github.Contributor) func() error {
	return func() error {
		options := &github.ListContributorsOptions{
			ListOptions: github.ListOptions{PerPage: 100},
		}
		var items []*github.Contributor
		for {
			data, resp, err := client.Repositories.ListContributors(c, repo.Owner, repo.RepoName, options)
			if err != nil {
				if resp != nil && resp.StatusCode == http.StatusNotFound {
					return &RepositoryNotFoundError{repo}
				}
				return err
			}
			items = append(items, data...)
			if resp.NextPage == 0 {
				break
			}
			options.Page = resp.NextPage
		}
		if out != nil && cap(out) > 0 {
			out <- items
		}
		return nil
	}
}
