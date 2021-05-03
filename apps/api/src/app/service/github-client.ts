import { Contributor, Repository } from '@lib/core';
import { Octokit } from '@octokit/rest';
import { injectable } from 'tsyringe';

@injectable()
export class GitHubClient {
  constructor(private readonly octokit: Octokit) {}

  async getAllContributors(repository: Repository): Promise<Contributor[]> {
    const { data } = await this.octokit.repos.listContributors({
      owner: repository.owner,
      repo: repository.repo,
      per_page: 100,
    });
    return data as Contributor[];
  }
}
