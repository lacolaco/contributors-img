import { Contributor, Repository } from '@lib/core';
import { injectable } from 'tsyringe';
import { CacheStorage } from '../service/cache-storage';
import { GitHubClient } from '../service/github-client';
import { runWithTracing } from '../utils/tracing';

function createCacheKey(repository: Repository, params: GetContributorsParams) {
  return `contributors-json-cache/${repository.owner}--${repository.repo}--${params.maxCount}.json`;
}

export interface GetContributorsParams {
  maxCount: number;
}

@injectable()
export class ContributorsRepository {
  constructor(private readonly githubClient: GitHubClient, private readonly cacheStorage: CacheStorage) {}

  async getAll(repository: Repository, params: GetContributorsParams): Promise<Contributor[]> {
    return runWithTracing('ContributorsRepository.getAllContributors', async () => {
      const cacheKey = createCacheKey(repository, params);
      const cached = await this.cacheStorage.restoreJSON<Contributor[]>(cacheKey);
      if (cached) {
        return cached;
      }

      const contributors = await this.githubClient.getContributors(repository, { maxCount: params.maxCount });
      this.cacheStorage.saveJSON(cacheKey, contributors);
      return contributors;
    });
  }
}
