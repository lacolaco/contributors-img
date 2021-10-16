import { Repository, RepositoryContributors } from '@lib/core';
import { injectable } from 'tsyringe';
import { CacheStorage } from '../service/cache-storage';
import { GitHubClient } from '../service/github-client';
import { runWithTracing } from '../utils/tracing';

function createCacheKey(repository: Repository, params: GetContributorsParams) {
  return `contributors-json-cache/v1.1/${repository.owner}--${repository.repo}--${params.maxCount}.json`;
}

export interface GetContributorsParams {
  maxCount: number;
}

@injectable()
export class ContributorsRepository {
  constructor(private readonly githubClient: GitHubClient, private readonly cacheStorage: CacheStorage) {}

  async getContributors(repository: Repository, params: GetContributorsParams): Promise<RepositoryContributors> {
    return runWithTracing('ContributorsRepository.getContributors', async () => {
      const cacheKey = createCacheKey(repository, params);
      const cached = await this.cacheStorage.restoreJSON<RepositoryContributors>(cacheKey);
      if (cached) {
        return cached;
      }

      const contributors = await this.githubClient.getContributors(repository, { maxCount: params.maxCount });
      await this.cacheStorage.saveJSON(cacheKey, contributors);
      return contributors;
    });
  }
}
