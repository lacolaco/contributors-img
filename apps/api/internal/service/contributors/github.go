package contributors

import (
	"context"
	"net/http"
	"strings"

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

	return buildRepositoryContributors(repositoryResult, contributorsResult), nil
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

func buildRepositoryContributors(rawRepository *github.Repository, rawContributors []*github.Contributor) *model.RepositoryContributors {
	contributors := make([]*model.Contributor, 0, len(rawContributors))
	for _, item := range rawContributors {
		// type is "User" or "Bot"
		if strings.ToLower(item.GetType()) == "bot" {
			continue
		}
		contributors = append(contributors, &model.Contributor{
			ID:            item.GetID(),
			Login:         item.GetLogin(),
			AvatarURL:     item.GetAvatarURL(),
			HTMLURL:       item.GetHTMLURL(),
			Contributions: item.GetContributions(),
		})
	}
	return &model.RepositoryContributors{
		Repository: &model.Repository{
			Owner:    rawRepository.Owner.GetLogin(),
			RepoName: rawRepository.GetName(),
		},
		StargazersCount: rawRepository.GetStargazersCount(),
		Contributors:    contributors,
	}
}
