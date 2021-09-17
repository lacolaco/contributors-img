import { Contributor, Repository } from '@lib/core';
import { injectable } from 'tsyringe';
import { ContributorsRepository, GetContributorsParams } from '../repository/contributors';
import { runWithTracing } from '../utils/tracing';

@injectable()
export class GetContributorsUsecase {
  constructor(private readonly contributorsRepository: ContributorsRepository) {}

  async execute(repoName: string, params: GetContributorsParams): Promise<Contributor[]> {
    return runWithTracing('GetContributorsUsecase.getContributors', async () => {
      const repository = Repository.fromString(repoName);

      return this.contributorsRepository.getAll(repository, params);
    });
  }
}
