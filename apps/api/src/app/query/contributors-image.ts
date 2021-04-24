import { Repository } from '@lib/core';
import { Readable } from 'stream';
import { injectable } from 'tsyringe';
import { ContributorsImageRepository } from '../repository/contributors-image';
import { UsageRepository } from '../repository/usage';
import { runWithTracing } from '../utils/tracing';
import { SupportedImageType } from '../utils/types';

@injectable()
export class ContributorsImageQuery {
  constructor(
    private readonly contributorsImageRepository: ContributorsImageRepository,
    private readonly usageRepository: UsageRepository,
  ) {}

  async getImage(
    repoName: string,
    options: { acceptWebp: boolean },
  ): Promise<{ fileStream: Readable; contentType: SupportedImageType }> {
    const repository = Repository.fromString(repoName);

    const { fileStream, contentType } = await runWithTracing('getImageFileStream', () =>
      this.contributorsImageRepository.getImageFileStream(repository, options),
    );

    runWithTracing('saveRepositoryUsage', () => this.usageRepository.saveRepositoryUsage(repository, new Date()));

    return { fileStream, contentType };
  }

  async getImage2(repoName: string): Promise<{ fileStream: Readable; contentType: SupportedImageType }> {
    const repository = Repository.fromString(repoName);

    const { fileStream, contentType } = await runWithTracing('getImageFileStream', () =>
      this.contributorsImageRepository.getImageFileStream2(repository),
    );

    runWithTracing('saveRepositoryUsage', () => this.usageRepository.saveRepositoryUsage(repository, new Date()));

    return { fileStream, contentType };
  }
}
