import { Repository } from '@lib/core';
import { injectable } from 'tsyringe';
import { ContributorsRepository } from '../repository/contributors';
import { ContributorsImageRepository } from '../repository/contributors-image';
import { UsageRepository } from '../repository/usage';

@injectable()
export class ContributorsImageQuery {
  constructor(
    private readonly contributorsRepository: ContributorsRepository,
    private readonly contributorsImageRepository: ContributorsImageRepository,
    private readonly usageRepository: UsageRepository,
  ) {}

  async getImage(repoName: string): Promise<Buffer> {
    const repository = Repository.fromString(repoName);

    const contributors = await this.contributorsRepository.getAllContributors(repository);
    const image = await this.contributorsImageRepository.getImage(repository, contributors);
    this.usageRepository.saveRepositoryUsage(repository, new Date());
    return image;
  }
}
