import { Repository } from '@lib/core';
import { injectable } from 'tsyringe';
import { ContributorsRepository } from '../repository/contributors';
import { ContributorsImageRepository } from '../repository/contributors-image';
import { UsageRepository } from '../repository/usage';
import { runWithTracing } from '../utils/tracing';

@injectable()
export class ContributorsImageQuery {
  constructor(
    private readonly contributorsRepository: ContributorsRepository,
    private readonly contributorsImageRepository: ContributorsImageRepository,
    private readonly usageRepository: UsageRepository,
  ) {}

  async getImage(repoName: string): Promise<Buffer> {
    const repository = Repository.fromString(repoName);

    const image = await runWithTracing('getImage', () => this.contributorsImageRepository.getImage(repository));

    runWithTracing('saveRepositoryUsage', () => this.usageRepository.saveRepositoryUsage(repository, new Date()));

    return image;
  }
}
