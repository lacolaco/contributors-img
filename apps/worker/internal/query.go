package app

import (
	"context"

	"contrib.rocks/libs/go/apiclient"
	"google.golang.org/api/iterator"
)

type RepositoryUsageRow struct {
	Owner        string `json:"owner"`
	Repository   string `json:"repository"`
	Days         int64  `json:"days"`
	Stargazers   int64  `json:"stargazers"`
	Contributors int64  `json:"contributors"`
}

func QueryFeaturedRepositories(ctx context.Context) ([]*RepositoryUsageRow, error) {
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
	rows := make([]*RepositoryUsageRow, 0, it.TotalRows)
	for {
		var row RepositoryUsageRow
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

type UsageStatsRow struct {
	Owners       int64 `json:"owners"`
	Repositories int64 `json:"repositories"`
}

func QueryUsageStats(ctx context.Context) (*UsageStatsRow, error) {
	bq := apiclient.NewBigQueryClient()
	q := bq.Query(`
  SELECT 
  	count(DISTINCT owner) owners,
  	count(DISTINCT repository) repositories,
  FROM
	` + "`contributors-img.repository_usage.weekly_repository_usage`" + `
  WHERE
    days >= 5`)
	it, err := q.Read(ctx)
	if err != nil {
		return nil, err
	}
	rows := make([]*UsageStatsRow, 0, it.TotalRows)
	for {
		var row UsageStatsRow
		err := it.Next(&row)
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		rows = append(rows, &row)
	}
	return rows[0], nil
}
