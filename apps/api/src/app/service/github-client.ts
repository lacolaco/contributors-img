import { Contributor, Repository } from '@lib/core';
import { Octokit } from '@octokit/rest';
import { concat, concatMap, filter, firstValueFrom, mapTo, of, range, toArray } from 'rxjs';
import { singleton } from 'tsyringe';
import { runWithTracing } from '../utils/tracing';

@singleton()
export class GitHubClient {
  constructor(private readonly octokit: Octokit) {}

  async getContributors(repository: Repository, { maxCount }: { maxCount: number }): Promise<Contributor[]> {
    return runWithTracing('GitHubClient.getContributors', async () => {
      const pages = Math.floor(maxCount / 100);
      const lastPageSize = maxCount % 100;
      const contributors = await firstValueFrom(
        concat(range(0, pages).pipe(mapTo(100)), of(lastPageSize)).pipe(
          filter((pageSize) => pageSize > 0),
          concatMap((pageSize, i) =>
            this.octokit.repos.listContributors({
              owner: repository.owner,
              repo: repository.repo,
              page: i + 1,
              per_page: pageSize,
            }),
          ),
          concatMap(({ data }) => of(...(data as Contributor[]))),
          toArray(),
        ),
      );

      return contributors;
    });
  }
}
