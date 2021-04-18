import { Contributor, Repository } from '@lib/core';
import { Octokit } from '@octokit/rest';
import { injectable } from 'tsyringe';

@injectable()
export class GitHubClient {
  constructor(private readonly octokit: Octokit) {}

  async getAllContributors(repository: Repository): Promise<Contributor[]> {
    const contributors = await this.octokit.paginate(this.octokit.repos.listContributors, {
      owner: repository.owner,
      repo: repository.repo,
      per_page: 100,
    });
    return contributors as Contributor[];
  }
}
