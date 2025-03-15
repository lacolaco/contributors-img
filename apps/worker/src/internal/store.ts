import { Firestore, CollectionReference } from '@google-cloud/firestore';
import { RepositoryUsageRow, UsageStatsRow } from './query';

export interface FeaturedRepository {
  repository: string;
  stargazers: number;
  contributors: number;
}

export interface FeaturedRepositoriesDocument {
  items: FeaturedRepository[];
  updatedAt: Date;
}

export interface UsageStats {
  owners: number;
  repositories: number;
}

export interface UsageStatsDocument {
  owners: number;
  repositories: number;
  updatedAt: Date;
}

function getEnvironmentCollection(appEnv: string): CollectionReference {
  const firestore = new Firestore();
  return firestore.collection(appEnv);
}

export async function saveFeaturedRepositories(
  appEnv: string,
  usageRows: RepositoryUsageRow[],
  updatedAt: Date,
): Promise<void> {
  const items: FeaturedRepository[] = usageRows.map((row) => ({
    repository: row.repository,
    stargazers: row.stargazers,
    contributors: row.contributors,
  }));

  const data: FeaturedRepositoriesDocument = {
    items,
    updatedAt,
  };

  await getEnvironmentCollection(appEnv).doc('featured_repositories').set(data);
}

export async function saveUsageStats(appEnv: string, usageRow: UsageStatsRow, updatedAt: Date): Promise<void> {
  const data: UsageStatsDocument = {
    owners: usageRow.owners,
    repositories: usageRow.repositories,
    updatedAt,
  };

  await getEnvironmentCollection(appEnv).doc('usage_stats').set(data);
}
