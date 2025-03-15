import { BigQuery } from '@google-cloud/bigquery';

export interface RepositoryUsageRow {
  owner: string;
  repository: string;
  days: number;
  stargazers: number;
  contributors: number;
}

export interface UsageStatsRow {
  owners: number;
  repositories: number;
}

export async function queryFeaturedRepositories(): Promise<RepositoryUsageRow[]> {
  const bq = new BigQuery();
  const query = `
    SELECT 
      owner,
      repository,
      days,
      usage.stargazers AS stargazers,
      usage.contributors AS contributors
    FROM
      \`contributors-img.repository_usage.weekly_repository_usage\`
    WHERE
      days >= 5
      AND usage.stargazers >= 1000
      AND usage.contributors >= 50
    LIMIT
      100`;

  const [rows] = await bq.query(query);
  return rows as RepositoryUsageRow[];
}

export async function queryUsageStats(): Promise<UsageStatsRow> {
  const bq = new BigQuery();
  const query = `
    SELECT 
      count(DISTINCT owner) AS owners,
      count(DISTINCT repository) AS repositories
    FROM
      \`contributors-img.repository_usage.weekly_repository_usage\`
    WHERE
      days >= 5`;

  const [rows] = await bq.query(query);
  return rows[0] as UsageStatsRow;
}
