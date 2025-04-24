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

type GitHubErrorType int

const (
	ErrorTypeUnknown GitHubErrorType = iota
	ErrorTypeTimeout
	ErrorTypeServer
	ErrorTypeConnectionRefused
	ErrorTypeNotFound
	ErrorTypeRateLimit
	ErrorTypeAbuseRateLimit
	ErrorTypeUnauthorized
	ErrorTypeForbidden
	ErrorTypeClientError
)

func getGitHubErrorType(err error, resp *github.Response) GitHubErrorType {
	if err == nil {
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			return ErrorTypeNotFound
		}
		return ErrorTypeUnknown
	}

	var netErr net.Error
	if errors.As(err, &netErr) && netErr.Timeout() {
		return ErrorTypeTimeout
	}

	var githubErr *github.ErrorResponse
	if errors.As(err, &githubErr) && githubErr.Response != nil {
		if githubErr.Response.StatusCode >= 500 && githubErr.Response.StatusCode < 600 {
			return ErrorTypeServer
		}
		if githubErr.Response.StatusCode == http.StatusNotFound {
			return ErrorTypeNotFound
		}
		if githubErr.Response.StatusCode == http.StatusUnauthorized {
			return ErrorTypeUnauthorized
		}
		if githubErr.Response.StatusCode == http.StatusForbidden {
			return ErrorTypeForbidden
		}
		if githubErr.Response.StatusCode >= 400 && githubErr.Response.StatusCode < 500 {
			return ErrorTypeClientError
		}
	}

	var rateLimitErr *github.RateLimitError
	if errors.As(err, &rateLimitErr) {
		return ErrorTypeRateLimit
	}

	var abuseErr *github.AbuseRateLimitError
	if errors.As(err, &abuseErr) {
		return ErrorTypeAbuseRateLimit
	}

	var opErr *net.OpError
	if errors.As(err, &opErr) {
		if opErr.Op == "dial" && strings.Contains(opErr.Error(), "connection refused") {
			return ErrorTypeConnectionRefused
		}
	}

	var urlErr *url.Error
	if errors.As(err, &urlErr) {
		if strings.Contains(urlErr.Error(), "connection refused") {
			return ErrorTypeConnectionRefused
		}
	}

	return ErrorTypeUnknown
}

func isRetryableError(err error) bool {
	errorType := getGitHubErrorType(err, nil)
	return errorType == ErrorTypeTimeout ||
		errorType == ErrorTypeServer ||
		errorType == ErrorTypeConnectionRefused ||
		errorType == ErrorTypeRateLimit ||
		errorType == ErrorTypeAbuseRateLimit
}

func isNotFoundError(err error, resp *github.Response) bool {
	return getGitHubErrorType(err, resp) == ErrorTypeNotFound
}

func getGitHubRetryOptions(ctx context.Context, repo *model.Repository) []retry.Option {
	return []retry.Option{
		retry.Attempts(3),
		retry.DelayType(retry.BackOffDelay),
		retry.RetryIf(func(err error) bool {
			return isRetryableError(err)
		}),
		retry.OnRetry(func(n uint, err error) {
			var rateLimitErr *github.RateLimitError
			if errors.As(err, &rateLimitErr) && rateLimitErr.Rate.Reset.Time.After(time.Now()) {
				waitTime := time.Until(rateLimitErr.Rate.Reset.Time)
				if waitTime > 0 {
					select {
					case <-time.After(waitTime):
					case <-ctx.Done():
						return
					}
				}
			}
		}),
		retry.Context(ctx),
	}
}

func fetchRepositoryContributors(client *github.Client, ctx context.Context, repo *model.Repository) (*model.RepositoryContributors, error) {
	ctx, span := tracing.Tracer().Start(ctx, "contributors.fetchRepositoryContributors")
	defer span.End()

	eg, groupCtx := errgroup.WithContext(ctx)

	var repository *github.Repository
	var contributors []*github.Contributor

	eg.Go(func() error {
		var resp *github.Response
		var err error

		err = retry.Do(
			func() error {
				var fetchErr error
				repository, resp, fetchErr = client.Repositories.Get(groupCtx, repo.Owner, repo.RepoName)
				return handleGitHubError(fetchErr, resp, repo)
			},
			getGitHubRetryOptions(groupCtx, repo)...,
		)

		return err
	})

	eg.Go(func() error {
		options := &github.ListContributorsOptions{
			Anon:        "true",
			ListOptions: github.ListOptions{PerPage: 100},
		}

		for {
			var data []*github.Contributor
			var resp *github.Response

			err := retry.Do(
				func() error {
					var fetchErr error
					data, resp, fetchErr = client.Repositories.ListContributors(groupCtx, repo.Owner, repo.RepoName, options)
					return handleGitHubError(fetchErr, resp, repo)
				},
				getGitHubRetryOptions(groupCtx, repo)...,
			)

			if err != nil {
				return err
			}

			contributors = append(contributors, data...)
			if resp.NextPage == 0 {
				break
			}
			options.Page = resp.NextPage
		}

		return nil
	})

	if err := eg.Wait(); err != nil {
		return nil, err
	}

	return buildRepositoryContributors(repository, contributors), nil
}

func handleGitHubError(err error, resp *github.Response, repo *model.Repository) error {
	if err != nil {
		if isNotFoundError(err, resp) {
			return retry.Unrecoverable(&RepositoryNotFoundError{repo})
		}
		return err
	}
	return nil
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
