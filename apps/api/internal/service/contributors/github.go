package contributors

import (
	"context"
	"crypto/md5"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"contrib.rocks/apps/api/go/model"
	"contrib.rocks/apps/api/internal/tracing"
	"github.com/avast/retry-go/v4"
	"github.com/google/go-github/v69/github"
	"golang.org/x/sync/errgroup"
)

func isTimeoutError(err error) bool {
	var netErr net.Error
	return errors.As(err, &netErr) && netErr.Timeout()
}

func isServerError(err error) bool {
	var githubErr *github.ErrorResponse
	if errors.As(err, &githubErr) && githubErr.Response != nil {
		return githubErr.Response.StatusCode >= 500 && githubErr.Response.StatusCode < 600
	}
	return false
}

func isConnectionRefusedError(err error) bool {
	var opErr *net.OpError
	if errors.As(err, &opErr) {
		return opErr.Op == "dial" && strings.Contains(opErr.Error(), "connection refused")
	}

	var urlErr *url.Error
	if errors.As(err, &urlErr) {
		return strings.Contains(urlErr.Error(), "connection refused")
	}

	return false
}

func getGitHubRetryOptions(ctx context.Context, repo *model.Repository) []retry.Option {
	return []retry.Option{
		retry.Attempts(3),
		retry.DelayType(retry.BackOffDelay),
		retry.RetryIf(func(err error) bool {
			var rateLimitErr *github.RateLimitError
			var abuseErr *github.AbuseRateLimitError

			return errors.As(err, &rateLimitErr) ||
				errors.As(err, &abuseErr) ||
				isServerError(err) ||
				isTimeoutError(err) ||
				isConnectionRefusedError(err)
		}),
		retry.OnRetry(func(n uint, err error) {
			var rateLimitErr *github.RateLimitError
			if errors.As(err, &rateLimitErr) && rateLimitErr.Rate.Reset.Time.After(time.Now()) {
				waitTime := time.Until(rateLimitErr.Rate.Reset.Time)
				if waitTime > 0 {
					time.Sleep(waitTime)
				}
			}
		}),
		retry.Context(ctx),
	}
}

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
		var data *github.Repository
		var resp *github.Response

		err := retry.Do(
			func() error {
				var err error
				data, resp, err = client.Repositories.Get(c, repo.Owner, repo.RepoName)
				if err != nil {
					if resp != nil && resp.StatusCode == http.StatusNotFound {
						return retry.Unrecoverable(&RepositoryNotFoundError{repo})
					}
					return err
				}
				return nil
			},
			getGitHubRetryOptions(c, repo)...,
		)

		if err != nil {
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
			Anon:        "true",
			ListOptions: github.ListOptions{PerPage: 100},
		}
		var items []*github.Contributor

		for {
			var data []*github.Contributor
			var resp *github.Response

			err := retry.Do(
				func() error {
					var err error
					data, resp, err = client.Repositories.ListContributors(c, repo.Owner, repo.RepoName, options)
					if err != nil {
						if resp != nil && resp.StatusCode == http.StatusNotFound {
							return retry.Unrecoverable(&RepositoryNotFoundError{repo})
						}
						return err
					}
					return nil
				},
				getGitHubRetryOptions(c, repo)...,
			)

			if err != nil {
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
