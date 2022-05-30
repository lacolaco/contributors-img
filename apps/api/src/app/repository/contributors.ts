import { Repository, RepositoryContributors } from '@lib/core';
import { injectable } from 'tsyringe';
import { CacheStorage } from '../service/cache-storage';
import { GitHubClient } from '../service/github-client';
import { runWithTracing } from '../utils/tracing';

function createCacheKey(repository: Repository) {
  return `contributors-json-cache/v1.2/${repository.owner}--${repository.repo}.json`;
}

@injectable()
export class ContributorsRepository {
  constructor(private readonly githubClient: GitHubClient, private readonly cacheStorage: CacheStorage) {}

  async getContributors(repository: Repository): Promise<RepositoryContributors> {
    return runWithTracing('ContributorsRepository.getContributors', async () => {
      const cacheKey = createCacheKey(repository);
      const cached = await this.cacheStorage.restoreJSON<RepositoryContributors>(cacheKey);
      if (cached) {
        return cached;
      }

      const contributors = await this.githubClient.getContributors(repository);
      await this.cacheStorage.saveJSON(cacheKey, contributors);
      return contributors;
    });
  }
}
