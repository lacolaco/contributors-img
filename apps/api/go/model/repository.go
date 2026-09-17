package model

import (
	"fmt"
	"regexp"
	"strings"
	"time"
)

// 'owner/repo'
type RepositoryString string

type Repository struct {
	Owner    string `json:"owner"`
	RepoName string `json:"repo"`
}

func (r RepositoryString) Object() *Repository {
	parts := strings.SplitN(string(r), "/", 2)
	return &Repository{Owner: parts[0], RepoName: parts[1]}
}

func (r Repository) String() string {
	return r.Owner + "/" + r.RepoName
}

// RepositoryNotFoundError is returned when a repository is not found
type RepositoryNotFoundError struct {
	Repository *Repository
}

func (e *RepositoryNotFoundError) Error() string {
	return "Repository not found: " + e.Repository.String()
}

// RateLimitedError is returned when the GitHub API rate limit is exhausted.
// RetryAfter is how long the caller should wait before the limit is expected
// to reset; callers use it to size a Retry-After / Cache-Control response.
type RateLimitedError struct {
	RetryAfter time.Duration
}

func (e *RateLimitedError) Error() string {
	return fmt.Sprintf("GitHub API rate limit exceeded, retry after %s", e.RetryAfter)
}

func ValidateRepositoryName(s string) error {
	if s == "" {
		return fmt.Errorf("repository name cannot be empty")
	}
	if match, err := regexp.MatchString(`^[\w\-._]+\/[\w\-._]+$`, s); !match || err != nil {
		return fmt.Errorf("invalid repository name: %s", s)
	}
	return nil
}
