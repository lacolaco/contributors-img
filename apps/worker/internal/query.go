package app

import (
	"context"

	"contrib.rocks/libs/goutils/apiclient"
	"google.golang.org/api/iterator"
)

type RepositoryUsage struct {
	Owner        string `json:"owner"`
	Repository   string `json:"repository"`
	Days         int64  `json:"days"`
	Stargazers   int64  `json:"stargazers"`
	Contributors int64  `json:"contributors"`
}

func QueryFeaturedRepositories(ctx context.Context) ([]*RepositoryUsage, error) {
	bq := apiclient.NewBigQueryClient()
	q := bq.Query(`
  SELECT 
	owner,
	repository,
	days,
	usage.stargazers AS stargazers,
	usage.contributors AS contributors
  FROM
	` + "`contributors-img.repository_usage.weekly_repository_usage`" + `
  WHERE
	days >= 5
	AND usage.stargazers >= 1000
	AND usage.contributors >= 50
  LIMIT
	100`)
	it, err := q.Read(ctx)
	if err != nil {
		return nil, err
	}
	rows := make([]*RepositoryUsage, 0, it.TotalRows)
	for {
		var row RepositoryUsage
		err := it.Next(&row)
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		rows = append(rows, &row)
	}
	return rows, nil
}
