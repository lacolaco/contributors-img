package api

import (
	"context"
	"errors"
	"net/http"
	"testing"
	"time"

	"contrib.rocks/apps/api/go/model"
	"github.com/google/go-github/v69/github"
)

// Test_Call_RateLimitError verifies that Call does not retry a rate limit
// error: it must invoke the call exactly once and return immediately,
// instead of waiting out the reset window like the old OnRetry hook did.
// Waiting here is exactly the behaviour that lets queued requests pile up
// and re-exhaust the rate limit the instant it resets.
func Test_Call_RateLimitError(t *testing.T) {
	rateLimitErr := &github.RateLimitError{
		Rate: github.Rate{Reset: github.Timestamp{Time: time.Now().Add(1 * time.Hour)}},
	}

	calls := 0
	start := time.Now()
	_, _, err := Call(context.Background(), func() (any, *github.Response, error) {
		calls++
		return nil, nil, rateLimitErr
	}, nil)
	elapsed := time.Since(start)

	if calls != 1 {
		t.Fatalf("expected call to be attempted once, got %d", calls)
	}
	if err == nil {
		t.Fatal("expected an error, got nil")
	}
	if elapsed > 1*time.Second {
		t.Fatalf("Call took %s, expected it to return immediately without waiting for the rate limit reset", elapsed)
	}
}

// Test_Call_AbuseRateLimitError mirrors Test_Call_RateLimitError for the
// secondary (abuse) rate limit error type.
func Test_Call_AbuseRateLimitError(t *testing.T) {
	retryAfter := 30 * time.Second
	abuseErr := &github.AbuseRateLimitError{RetryAfter: &retryAfter}

	calls := 0
	_, _, err := Call(context.Background(), func() (any, *github.Response, error) {
		calls++
		return nil, nil, abuseErr
	}, nil)

	if calls != 1 {
		t.Fatalf("expected call to be attempted once, got %d", calls)
	}
	if err == nil {
		t.Fatal("expected an error, got nil")
	}
}

// Test_Call_ServerError_StillRetries confirms the existing backoff-retry
// behaviour for 5xx errors is unchanged by removing the rate limit branches.
func Test_Call_ServerError_StillRetries(t *testing.T) {
	serverErr := &github.ErrorResponse{Response: &http.Response{StatusCode: http.StatusInternalServerError}}

	calls := 0
	_, _, err := Call(context.Background(), func() (any, *github.Response, error) {
		calls++
		return nil, nil, serverErr
	}, nil)

	if calls != 3 {
		t.Fatalf("expected 3 attempts (retry.Attempts(3)), got %d", calls)
	}
	if err == nil {
		t.Fatal("expected an error, got nil")
	}
}

func Test_HandleRateLimitError(t *testing.T) {
	t.Run("nil error passes through", func(t *testing.T) {
		if err := HandleRateLimitError(nil); err != nil {
			t.Fatalf("expected nil, got %v", err)
		}
	})

	t.Run("unrelated error passes through unchanged", func(t *testing.T) {
		original := context.DeadlineExceeded
		if err := HandleRateLimitError(original); err != original {
			t.Fatalf("expected original error to pass through, got %v", err)
		}
	})

	t.Run("RateLimitError with future reset uses time until reset", func(t *testing.T) {
		reset := time.Now().Add(2 * time.Minute)
		err := HandleRateLimitError(&github.RateLimitError{Rate: github.Rate{Reset: github.Timestamp{Time: reset}}})

		var rateLimited *model.RateLimitedError
		if !errors.As(err, &rateLimited) {
			t.Fatalf("expected *model.RateLimitedError, got %T: %v", err, err)
		}
		if rateLimited.RetryAfter <= 0 || rateLimited.RetryAfter > 2*time.Minute {
			t.Fatalf("expected RetryAfter close to 2m, got %s", rateLimited.RetryAfter)
		}
	})

	t.Run("RateLimitError with past reset falls back to the default", func(t *testing.T) {
		reset := time.Now().Add(-1 * time.Minute)
		err := HandleRateLimitError(&github.RateLimitError{Rate: github.Rate{Reset: github.Timestamp{Time: reset}}})

		var rateLimited *model.RateLimitedError
		if !errors.As(err, &rateLimited) {
			t.Fatalf("expected *model.RateLimitedError, got %T: %v", err, err)
		}
		if rateLimited.RetryAfter != defaultRateLimitRetryAfter {
			t.Fatalf("expected default RetryAfter %s, got %s", defaultRateLimitRetryAfter, rateLimited.RetryAfter)
		}
	})

	t.Run("AbuseRateLimitError without RetryAfter uses default", func(t *testing.T) {
		err := HandleRateLimitError(&github.AbuseRateLimitError{})

		var rateLimited *model.RateLimitedError
		if !errors.As(err, &rateLimited) {
			t.Fatalf("expected *model.RateLimitedError, got %T: %v", err, err)
		}
		if rateLimited.RetryAfter != defaultRateLimitRetryAfter {
			t.Fatalf("expected default RetryAfter %s, got %s", defaultRateLimitRetryAfter, rateLimited.RetryAfter)
		}
	})

	t.Run("AbuseRateLimitError with RetryAfter uses it", func(t *testing.T) {
		retryAfter := 45 * time.Second
		err := HandleRateLimitError(&github.AbuseRateLimitError{RetryAfter: &retryAfter})

		var rateLimited *model.RateLimitedError
		if !errors.As(err, &rateLimited) {
			t.Fatalf("expected *model.RateLimitedError, got %T: %v", err, err)
		}
		if rateLimited.RetryAfter != retryAfter {
			t.Fatalf("expected RetryAfter %s, got %s", retryAfter, rateLimited.RetryAfter)
		}
	})
}
