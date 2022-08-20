package model

type RepositoryContributors struct {
	*Repository
	StargazersCount int            `json:"stargazersCount"`
	Contributors    []*Contributor `json:"data"`
}

type Contributor struct {
	ID            int64  `json:"id"`
	Login         string `json:"login"`
	AvatarURL     string `json:"avatar_url"`
	HTMLURL       string `json:"html_url"`
	Contributions int    `json:"contributions"`
}
