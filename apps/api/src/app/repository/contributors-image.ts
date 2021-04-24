import { Repository } from '@lib/core';
import { Readable } from 'stream';
import { injectable } from 'tsyringe';
import { CacheStorage } from '../service/cache-storage';
import { ContributorsImageSvgRenderer } from '../service/contributors-image-renderer';
import { runWithTracing } from '../utils/tracing';
import { ContributorsRepository } from './contributors';

function createCacheKey(repository: Repository, ext: string) {
  return `image-cache/${repository.owner}--${repository.repo}.${ext}`;
}

@injectable()
export class ContributorsImageRepository {
  constructor(
    private readonly cacheStorage: CacheStorage,
    private readonly contributorsRepository: ContributorsRepository,
    private readonly renderer: ContributorsImageSvgRenderer,
  ) {}

  async getImageFileStream(repository: Repository): Promise<{ fileStream: Readable; contentType: string }> {
    const cacheKey = createCacheKey(repository, 'svg');
    const cached = await runWithTracing('restoreFromCache', async () => {
      return this.cacheStorage
        .restoreFileStream(cacheKey)
        .then((fileStream) => (fileStream ? { fileStream, contentType: 'image/svg+xml' } : null));
    });
    if (cached) {
      return cached;
    }

    const contributors = await runWithTracing('getAllContributors', () =>
      this.contributorsRepository.getAllContributors(repository),
    );

    const svgImage = await runWithTracing('renderImage', async () =>
      this.renderer.renderSimpleAvatarTable(contributors),
    );

    runWithTracing('saveCache', async () => {
      await Promise.all([this.cacheStorage.saveFile(cacheKey, svgImage, 'image/svg+xml')]);
    });

    return {
      fileStream: Readable.from(svgImage),
      contentType: 'image/svg+xml',
    };
  }
}
