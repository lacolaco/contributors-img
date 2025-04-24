// Package api provides utilities for GitHub API operations
package api

import (
	"context"
	"errors"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"contrib.rocks/apps/api/go/model"
	"github.com/avast/retry-go/v4"
	"github.com/google/go-github/v69/github"
)

// ErrorType represents the categorized GitHub error types
type ErrorType int

const (
	ErrorTypeUnknown ErrorType = iota
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

// errorTypeMapping defines a mapping between error conditions and error types
var errorTypeMapping = []struct {
	check     func(error) bool
	errorType ErrorType
}{
	{
		check: func(err error) bool {
			var netErr net.Error
			return errors.As(err, &netErr) && netErr.Timeout()
		},
		errorType: ErrorTypeTimeout,
	},
	{
		check: func(err error) bool {
			var rateLimitErr *github.RateLimitError
			return errors.As(err, &rateLimitErr)
		},
		errorType: ErrorTypeRateLimit,
	},
	{
		check: func(err error) bool {
			var abuseErr *github.AbuseRateLimitError
			return errors.As(err, &abuseErr)
		},
		errorType: ErrorTypeAbuseRateLimit,
	},
	{
		check: func(err error) bool {
			var opErr *net.OpError
			return errors.As(err, &opErr) && opErr.Op == "dial" &&
				strings.Contains(opErr.Error(), "connection refused")
		},
		errorType: ErrorTypeConnectionRefused,
	},
	{
		check: func(err error) bool {
			var urlErr *url.Error
			return errors.As(err, &urlErr) &&
				strings.Contains(urlErr.Error(), "connection refused")
		},
		errorType: ErrorTypeConnectionRefused,
	},
}

// statusCodeMapping maps HTTP status codes to error types
var statusCodeMapping = map[int]ErrorType{
	http.StatusNotFound:     ErrorTypeNotFound,
	http.StatusUnauthorized: ErrorTypeUnauthorized,
	http.StatusForbidden:    ErrorTypeForbidden,
}

// GetErrorType determines the type of GitHub error
func GetErrorType(err error, resp *github.Response) ErrorType {
	if err == nil {
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			return ErrorTypeNotFound
		}
		return ErrorTypeUnknown
	}

	// Check specific error types
	for _, mapping := range errorTypeMapping {
		if mapping.check(err) {
			return mapping.errorType
		}
	}

	// Check GitHub API response status codes
	var githubErr *github.ErrorResponse
	if errors.As(err, &githubErr) && githubErr.Response != nil {
		statusCode := githubErr.Response.StatusCode

		if errorType, ok := statusCodeMapping[statusCode]; ok {
			return errorType
		}

		if statusCode >= 500 {
			return ErrorTypeServer
		}
		if statusCode >= 400 {
			return ErrorTypeClientError
		}
	}

	return ErrorTypeUnknown
}

// IsRetryableError checks if an error is retryable
func IsRetryableError(err error) bool {
	errorType := GetErrorType(err, nil)
	return errorType == ErrorTypeTimeout ||
		errorType == ErrorTypeServer ||
		errorType == ErrorTypeConnectionRefused ||
		errorType == ErrorTypeRateLimit ||
		errorType == ErrorTypeAbuseRateLimit
}

// IsNotFoundError checks if an error is a not found error
func IsNotFoundError(err error, resp *github.Response) bool {
	return GetErrorType(err, resp) == ErrorTypeNotFound
}

// GetRetryOptions returns standard retry options for GitHub API calls
func GetRetryOptions(ctx context.Context) []retry.Option {
	return []retry.Option{
		retry.Attempts(3),
		retry.DelayType(retry.BackOffDelay),
		retry.RetryIf(func(err error) bool {
			return IsRetryableError(err)
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

// Call executes a GitHub API call with retry logic
func Call[T any](
	ctx context.Context,
	call func() (T, *github.Response, error),
	errorHandler func(error, *github.Response) error,
) (T, *github.Response, error) {
	var result T
	var resp *github.Response

	err := retry.Do(
		func() error {
			var callErr error
			result, resp, callErr = call()
			if errorHandler != nil {
				return errorHandler(callErr, resp)
			}
			return callErr
		},
		GetRetryOptions(ctx)...,
	)

	if err != nil {
		return result, resp, err
	}

	return result, resp, nil
}

// HandleRepositoryNotFoundError converts not found errors to a RepositoryNotFoundError
func HandleRepositoryNotFoundError(err error, resp *github.Response, repo *model.Repository) error {
	if err != nil && IsNotFoundError(err, resp) {
		return retry.Unrecoverable(&model.RepositoryNotFoundError{Repository: repo})
	}
	return err
}
