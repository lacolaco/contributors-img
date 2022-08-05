package model

import "strings"

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

func (r Repository) String() RepositoryString {
	return RepositoryString(r.Owner + "/" + r.RepoName)
}
