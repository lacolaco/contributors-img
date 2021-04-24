import { Repository } from '@lib/core';
import { Readable } from 'stream';
import { injectable } from 'tsyringe';
import { ContributorsImageRepository } from '../repository/contributors-image';
import { UsageRepository } from '../repository/usage';
import { runWithTracing } from '../utils/tracing';

@injectable()
export class ContributorsImageQuery {
  constructor(
    private readonly contributorsImageRepository: ContributorsImageRepository,
    private readonly usageRepository: UsageRepository,
  ) {}

  async getImage(repoName: string): Promise<{ fileStream: Readable; contentType: string }> {
    const repository = Repository.fromString(repoName);

    const { fileStream, contentType } = await runWithTracing('getImageFileStream', () =>
      this.contributorsImageRepository.getImageFileStream(repository),
    );

    runWithTracing('saveRepositoryUsage', () => this.usageRepository.saveRepositoryUsage(repository, new Date()));

    return { fileStream, contentType };
  }
}
