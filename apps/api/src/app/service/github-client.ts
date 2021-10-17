import { Contributor, Repository, RepositoryContributors } from '@lib/core';
import { Octokit } from '@octokit/rest';
import { concatMap, firstValueFrom, forkJoin, from, map, Observable, of, toArray } from 'rxjs';
import { singleton } from 'tsyringe';
import { createPager } from '../utils/paging';
import { runWithTracing } from '../utils/tracing';

@singleton()
export class GitHubClient {
  constructor(private readonly octokit: Octokit) {}

  async getContributors(repository: Repository, { maxCount }: { maxCount: number }): Promise<RepositoryContributors> {
    return runWithTracing('GitHubClient.getContributors', async () => {
      const data$ = forkJoin([this.fetchRepositoryMeta(repository), this.fetchContributors(repository, maxCount)]).pipe(
        map(([{ stargazersCount }, contributors]) => ({
          ...repository,
          stargazersCount,
          data: contributors,
        })),
      );

      return firstValueFrom(data$);
    });
  }
  private fetchRepositoryMeta({ owner, repo }: Repository): Observable<{ stargazersCount: number }> {
    return from(this.octokit.repos.get({ owner, repo })).pipe(
      map(({ data }) => ({
        stargazersCount: data.stargazers_count,
      })),
    );
  }

  private fetchContributors({ owner, repo }: Repository, maxCount: number): Observable<Contributor[]> {
    const pages = createPager(100)(maxCount);

    return from(pages).pipe(
      concatMap((pageSize, i) =>
        this.octokit.repos.listContributors({
          owner,
          repo,
          page: i + 1,
          per_page: pageSize,
        }),
      ),
      concatMap(({ data }) => of(...(data as Contributor[]))),
      toArray(),
    );
  }
}
