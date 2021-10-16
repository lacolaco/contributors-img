import { Firestore } from '@google-cloud/firestore';
import { Logging } from '@google-cloud/logging';
import { Repository, RepositoryContributors } from '@lib/core';
import { injectable } from 'tsyringe';
import { environment } from '../../environments/environment';
import { runWithTracing } from '../utils/tracing';

@injectable()
export class UsageCollector {
  constructor(private readonly firestore: Firestore, private readonly logging: Logging) {}

  async collectUsage({ owner, repo, stargazersCount, data }: RepositoryContributors, timestamp = Date.now()) {
    return runWithTracing('UsageCollector.collectUsage', async () => {
      const log = this.logging.log('repository-usage');
      const entry = log.entry(
        {
          labels: {
            environment: environment.environmentName,
          },
        },
        {
          repository: `${owner}/${repo}`,
          stargazers: stargazersCount,
          contributors: data.length,
          timestamp: timestamp.toString(),
        },
      );
      return await log.write(entry);
    });
  }

  async saveRepositoryUsage(repository: Repository, timestamp: Date) {
    return runWithTracing('UsageCollector.saveRepositoryUsage', async () => {
      try {
        await this.firestore
          .collection(`${environment.firestoreRootCollectionName}/usage/repositories`)
          .doc(`${repository.owner}--${repository.repo}`)
          .set({
            name: repository.toString(),
            timestamp,
          });
      } catch (error) {
        console.error(error);
      }
    });
  }
}
