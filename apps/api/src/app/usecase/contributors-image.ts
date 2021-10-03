import { Repository } from '@lib/core';
import { Readable } from 'stream';
import { injectable } from 'tsyringe';
import { ContributorsRepository } from '../repository/contributors';
import { ContributorsImageRepository } from '../repository/contributors-image';
import { ContributorsImageRenderer } from '../service/contributors-image-renderer';
import { UsageCollector as UsageCollector } from '../service/usage-collector';
import { runWithTracing } from '../utils/tracing';
import { FileStream } from '../utils/types';
import { defaultContributorsMaxCount } from './constants';

@injectable()
export class GetContributorsImageUsecase {
  constructor(
    private readonly contributorsImageRepository: ContributorsImageRepository,
    private readonly contributorsRepository: ContributorsRepository,
    private readonly usageCollector: UsageCollector,
    private readonly renderer: ContributorsImageRenderer,
  ) {}

  async execute({
    repository,
    isGitHubRequest,
    preview = false,
    maxCount = null,
  }: {
    repository: Repository;
    preview: boolean;
    isGitHubRequest: boolean;
    maxCount: number | null;
  }): Promise<FileStream> {
    return runWithTracing('GetContributorsImageUsecase.getImage', async () => {
      maxCount = maxCount ?? defaultContributorsMaxCount;

      // restore cached image if exists
      const image = await this.contributorsImageRepository.loadImage(repository, { maxCount });
      if (image.data) {
        return image;
      }

      // get contributors
      const contributors = await this.contributorsRepository.getContributors(repository, { maxCount });
      // render image
      const { data, contentType } = await this.renderer.render(contributors);
      // save image to cache
      await image.save(data, contentType);

      if (isGitHubRequest) {
        this.usageCollector.collectUsage(repository);
      }
      if (!preview) {
        this.usageCollector.saveRepositoryUsage(repository, new Date());
      }

      return {
        data: Readable.from(data),
        contentType,
      };
    });
  }
}
