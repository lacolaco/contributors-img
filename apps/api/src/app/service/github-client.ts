import { Contributor, Repository, RepositoryContributors } from '@lib/core';
import { Octokit } from '@octokit/rest';
import { firstValueFrom, forkJoin, from, map, Observable } from 'rxjs';
import { singleton } from 'tsyringe';
import { runWithTracing } from '../utils/tracing';

@singleton()
export class GitHubClient {
  constructor(private readonly octokit: Octokit) {}

  async getContributors(repository: Repository): Promise<RepositoryContributors> {
    return runWithTracing('GitHubClient.getContributors', async () => {
      const data$ = forkJoin([this.fetchRepositoryMeta(repository), this.fetchContributors(repository)]).pipe(
        map(([meta, contributors]) => ({
          ...meta,
          data: contributors,
        })),
      );

      return firstValueFrom(data$);
    });
  }
  private fetchRepositoryMeta({
    owner,
    repo,
  }: Repository): Observable<{ owner: string; repo: string; stargazersCount: number }> {
    return from(this.octokit.repos.get({ owner, repo })).pipe(
      map(({ data }) => ({
        owner: data.owner.login,
        repo: data.name,
        stargazersCount: data.stargazers_count,
      })),
    );
  }

  private async fetchContributors({ owner, repo }: Repository): Promise<Contributor[]> {
    return (await this.octokit.paginate(this.octokit.repos.listContributors, {
      owner,
      repo,
      per_page: 100,
    })) as Contributor[];
  }
}
