import { Repository } from '@lib/core';
import { container } from 'tsyringe';
import { CacheStorage } from '../service/cache-storage';
import { GitHubClient } from '../service/github-client';
import { ContributorsRepository } from './contributors';

describe('ContributorsRepository', () => {
  let repository: ContributorsRepository;

  beforeEach(() => {
    container.clearInstances();
    container.registerInstance(CacheStorage, new CacheStorage(null));
    repository = container.resolve(ContributorsRepository);
  });

  describe('getAll', () => {
    test('it gets contributors from github api', async () => {
      const repo = new Repository('angular', 'angular');
      const githubClient = container.resolve(GitHubClient);
      jest.spyOn(githubClient, 'getContributors').mockResolvedValue({ ...repo, data: [], stargazersCount: 0 });

      await repository.getContributors(repo);

      expect(githubClient.getContributors).toHaveBeenCalled();
    });
  });
});
