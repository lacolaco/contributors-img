import { Contributor, Repository } from '@lib/core';
import { injectable } from 'tsyringe';
import { ContributorsRepository } from '../repository/contributors';
import { runWithTracing } from '../utils/tracing';

@injectable()
export class ContributorsQuery {
  constructor(private readonly contributorsRepository: ContributorsRepository) {}

  async getContributors(repoName: string): Promise<Contributor[]> {
    const repository = Repository.fromString(repoName);

    const contributors = await runWithTracing('getAllContributors', () =>
      this.contributorsRepository.getAllContributors(repository),
    );
    return contributors;
  }
}
