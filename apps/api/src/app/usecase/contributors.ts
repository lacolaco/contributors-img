import { Contributor, Repository } from '@lib/core';
import { injectable } from 'tsyringe';
import { ContributorsRepository } from '../repository/contributors';
import { runWithTracing } from '../utils/tracing';
import { defaultContributorsMaxCount } from './constants';

@injectable()
export class GetContributorsUsecase {
  constructor(private readonly contributorsRepository: ContributorsRepository) {}

  async execute(repository: Repository, maxCount: number | null): Promise<Contributor[]> {
    return runWithTracing('GetContributorsUsecase.getContributors', async () => {
      maxCount = maxCount ?? defaultContributorsMaxCount;

      return this.contributorsRepository.getContributors(repository, { maxCount });
    });
  }
}
