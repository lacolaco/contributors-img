package contributors

import (
	"context"
	"crypto/md5"
	"fmt"
	"strings"

	"contrib.rocks/apps/api/go/model"
	"contrib.rocks/apps/api/internal/github/api"
	"contrib.rocks/apps/api/internal/tracing"
	"github.com/google/go-github/v69/github"
	"golang.org/x/sync/errgroup"
)

func fetchRepositoryContributors(client *github.Client, ctx context.Context, repo *model.Repository) (*model.RepositoryContributors, error) {
	ctx, span := tracing.Tracer().Start(ctx, "contributors.fetchRepositoryContributors")
	defer span.End()

	eg, groupCtx := errgroup.WithContext(ctx)

	var repository *github.Repository
	var contributors []*github.Contributor

	// 並列でリポジトリ情報とコントリビューター情報を取得
	eg.Go(func() error {
		var err error
		repository, err = fetchRepository(client, groupCtx, repo)
		return err
	})

	eg.Go(func() error {
		var err error
		contributors, err = fetchAllContributors(client, groupCtx, repo)
		return err
	})

	if err := eg.Wait(); err != nil {
		return nil, err
	}

	return buildRepositoryContributors(repository, contributors), nil
}

// リポジトリ情報を取得
func fetchRepository(client *github.Client, ctx context.Context, repo *model.Repository) (*github.Repository, error) {
	errorHandler := func(err error, resp *github.Response) error {
		return api.HandleRepositoryNotFoundError(err, resp, repo)
	}

	repository, _, err := api.Call(ctx, func() (*github.Repository, *github.Response, error) {
		return client.Repositories.Get(ctx, repo.Owner, repo.RepoName)
	}, errorHandler)

	return repository, err
}

// すべてのコントリビューター情報をページングしながら取得
func fetchAllContributors(client *github.Client, ctx context.Context, repo *model.Repository) ([]*github.Contributor, error) {
	var allContributors []*github.Contributor
	options := &github.ListContributorsOptions{
		Anon:        "true",
		ListOptions: github.ListOptions{PerPage: 100},
	}

	for {
		errorHandler := func(err error, resp *github.Response) error {
			return api.HandleRepositoryNotFoundError(err, resp, repo)
		}

		contributors, resp, err := api.Call(ctx, func() ([]*github.Contributor, *github.Response, error) {
			return client.Repositories.ListContributors(ctx, repo.Owner, repo.RepoName, options)
		}, errorHandler)

		if err != nil {
			return nil, err
		}

		allContributors = append(allContributors, contributors...)

		// レスポンスからページング情報を取得
		if resp == nil || resp.NextPage == 0 {
			break
		}
		options.Page = resp.NextPage
	}

	return allContributors, nil
}

func buildRepositoryContributors(rawRepository *github.Repository, rawContributors []*github.Contributor) *model.RepositoryContributors {
	contributors := make([]*model.Contributor, 0, len(rawContributors))
	for _, item := range rawContributors {
		switch strings.ToLower(item.GetType()) {
		case "bot":
			continue
		case "anonymous":
			contributors = append(contributors, &model.Contributor{
				Login:         item.GetName(),
				AvatarURL:     fmt.Sprintf("https://www.gravatar.com/avatar/%x?d=mp", md5.Sum([]byte(item.GetEmail()))),
				Contributions: item.GetContributions(),
			})
			continue
		default:
			contributors = append(contributors, &model.Contributor{
				ID:            item.GetID(),
				Login:         item.GetLogin(),
				AvatarURL:     item.GetAvatarURL(),
				HTMLURL:       item.GetHTMLURL(),
				Contributions: item.GetContributions(),
			})
		}
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
