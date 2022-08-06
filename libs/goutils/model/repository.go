package model

import (
	"fmt"
	"regexp"
	"strings"
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

func ValidateRepositoryName(s string) error {
	if s == "" {
		return fmt.Errorf("repository name cannot be empty")
	}
	if match, err := regexp.MatchString(`^[\w\-._]+\/[\w\-._]+$`, s); !match || err != nil {
		return fmt.Errorf("invalid repository name: %s", s)
	}
	return nil
}
