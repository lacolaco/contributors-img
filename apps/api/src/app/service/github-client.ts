import { Contributor, Repository } from '@lib/core';
import { Octokit } from '@octokit/rest';
import { concat, of, range } from 'rxjs';
import { concatMap, filter, mapTo, toArray } from 'rxjs/operators';
import { singleton } from 'tsyringe';

@singleton()
export class GitHubClient {
  constructor(private readonly octokit: Octokit) {}

  async getContributors(repository: Repository, { maxCount }: { maxCount: number }): Promise<Contributor[]> {
    const pages = Math.floor(maxCount / 100);
    const lastPageSize = maxCount % 100;
    const contributors = await concat(range(0, pages).pipe(mapTo(100)), of(lastPageSize))
      .pipe(
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
      )
      .toPromise();

    return contributors;
  }
}
