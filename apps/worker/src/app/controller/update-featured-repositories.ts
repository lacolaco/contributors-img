import { BigQuery } from '@google-cloud/bigquery';
import { Firestore } from '@google-cloud/firestore';
import { Request, Response } from 'express';
import { injectable } from 'tsyringe';
import { environment } from '../../environments/environment';

type RepositoryUsageRow = Readonly<{
  owner: string;
  days: number;
  usage: {
    repository: string;
    stargazers: number;
    contributors: number;
  };
}>;

async function queryRepositoryUsage({ minStars = 1000, limit = 50 }: { minStars?: number; limit?: number }) {
  const bq = new BigQuery();
  const query = `
SELECT
  owner,
  COUNT(DISTINCT _date) AS days,
  ARRAY_AGG(STRUCT(repository, stargazers, contributors) ORDER BY contributors DESC, stargazers DESC)[OFFSET(0)] as usage,
FROM (
  SELECT
    jsonPayload.repository AS repository,
    SPLIT(jsonPayload.repository, '/')[OFFSET(0)] AS owner,
    DATE(TIMESTAMP_MILLIS(CAST(jsonPayload.timestamp AS INT64))) AS _date,
    jsonPayload.stargazers AS stargazers,
    jsonPayload.contributors AS contributors
  FROM
    \`contributors-img.repository_usage.repository_usage_*\`
  WHERE
    labels.environment = @environment
    AND _TABLE_SUFFIX BETWEEN FORMAT_DATE('%Y%m%d', DATE_SUB(CURRENT_DATE(), INTERVAL 7 day))
    AND FORMAT_DATE('%Y%m%d', CURRENT_DATE()))
GROUP BY
  owner
HAVING
  days >= 6
  AND usage.stargazers >= @minStars
  AND usage.contributors >= 50
ORDER BY
  usage.stargazers DESC,
  usage.contributors DESC
LIMIT
  @limit`.trim();
  const [rows] = await bq.query(
    {
      query,
      params: {
        environment: environment.environment,
        minStars,
        limit,
      },
    },
    {},
  );

  return rows as RepositoryUsageRow[];
}

async function saveFeaturedRepositories(featuredRepositories: RepositoryUsageRow[], updatedAt: Date) {
  const firestore = new Firestore();
  try {
    await firestore
      .collection(`${environment.environment}`)
      .doc('featured_repositories')
      .set({
        items: featuredRepositories.map((item) => item.usage),
        updatedAt,
      });
  } catch (error) {
    console.error(error);
  }
}

@injectable()
export class UpdateFeaturedRepositoriesController {
  async onRequest(req: Request, res: Response) {
    try {
      const rows = await queryRepositoryUsage({});
      rows.forEach((row) => console.log(JSON.stringify(row)));
      await saveFeaturedRepositories(rows, new Date());
      res.status(200).send('OK');
    } catch (error) {
      res.status(500).send(error);
    }
  }
}
