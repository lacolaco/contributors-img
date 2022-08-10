package app

import (
	"context"

	"cloud.google.com/go/bigquery"
	"contrib.rocks/libs/goutils/apiclient"
	"contrib.rocks/libs/goutils/env"
	"google.golang.org/api/iterator"
)

func QueryFeaturedRepositories(ctx context.Context, appEnv env.Environment) ([]FeaturedRepository, error) {
	bq := apiclient.NewBigQueryClient()
	q := bq.Query(`
	SELECT
	  owner,
	  COUNT(DISTINCT _date) AS days,
	  ARRAY_AGG(STRUCT(repository, stargazers, contributors) ORDER BY contributors DESC, stargazers DESC)[OFFSET(0)] as usage,
	FROM (
	  SELECT
		jsonPayload.owner AS owner,
		jsonPayload.repository AS repository,
		DATE(timestamp) AS _date,
		CAST(jsonPayload.stargazers as INT64) AS stargazers,
		CAST(jsonPayload.contributors as INT64) AS contributors
	  FROM ` + "`contributors-img.repository_usage.repository_usage_*`" + `
	  WHERE
		labels.environment = @environment
		AND labels.via = "github"
		AND _TABLE_SUFFIX BETWEEN FORMAT_DATE('%Y%m%d', DATE_SUB(CURRENT_DATE(), INTERVAL @dayInterval day))
		AND FORMAT_DATE('%Y%m%d', CURRENT_DATE()))
	GROUP BY
	  owner
	HAVING
	  days >= @minDays
	  AND usage.stargazers >= @minStars
	  AND usage.contributors >= 50
	ORDER BY
	  usage.stargazers DESC,
	  usage.contributors DESC
	LIMIT
	  @limit`)
	q.Parameters = []bigquery.QueryParameter{
		{Name: "environment", Value: string(appEnv)},
		{Name: "dayInterval", Value: int(6)},
		{Name: "minDays", Value: int(5)},
		{Name: "minStars", Value: int(1000)},
		{Name: "limit", Value: int(50)},
	}
	it, err := q.Read(ctx)
	if err != nil {
		return nil, err
	}
	rows := make([]FeaturedRepository, 0, it.TotalRows)
	for {
		var row FeaturedRepository
		err := it.Next(&row)
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		rows = append(rows, row)
	}
	return rows, nil
}
