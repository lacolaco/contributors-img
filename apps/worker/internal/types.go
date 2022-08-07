package app

type RepositoryUsage struct {
	Repository   string  `json:"repository" firestore:"repository"`
	Stargazers   float32 `json:"stargazers" firestore:"stargazers"`
	Contributors float32 `json:"contributors" firestore:"contributors"`
}

type FeaturedRepository struct {
	Owner string          `json:"owner" firestore:"owner"`
	Days  int64           `json:"days" firestore:"days"`
	Usage RepositoryUsage `json:"usage" firestore:"usage"`
}
