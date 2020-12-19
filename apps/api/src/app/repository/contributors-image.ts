import { Repository } from '@lib/core';
import { Readable } from 'stream';
import * as sharp from 'sharp';
import { injectable } from 'tsyringe';
import { CacheStorage } from '../service/cache-storage';
import { ContributorsImageRenderer } from '../service/contributors-image-renderer';
import { runWithTracing } from '../utils/tracing';
import { SupportedImageType } from '../utils/types';
import { ContributorsRepository } from './contributors';

function createCacheKey(repository: Repository, ext: string) {
  return `image-cache/${repository.owner}--${repository.repo}.${ext}`;
}

@injectable()
export class ContributorsImageRepository {
  constructor(
    private readonly cacheStorage: CacheStorage,
    private readonly contributorsRepository: ContributorsRepository,
    private readonly renderer: ContributorsImageRenderer,
  ) {}

  async getImageFileStream(
    repository: Repository,
    options: { acceptWebp: boolean },
  ): Promise<{ fileStream: Readable; contentType: SupportedImageType }> {
    const cached = await runWithTracing('restoreFromCache', async () => {
      const pngKey = createCacheKey(repository, 'png');
      const webpKey = createCacheKey(repository, 'webp');
      if (options.acceptWebp) {
        const caches = await Promise.all([
          this.cacheStorage
            .restoreFileStream(webpKey)
            .then((fileStream) =>
              fileStream ? { fileStream, contentType: 'image/webp' as SupportedImageType } : null,
            ),
          this.cacheStorage
            .restoreFileStream(pngKey)
            .then((fileStream) => (fileStream ? { fileStream, contentType: 'image/png' as SupportedImageType } : null)),
        ]);
        return caches.find((cache) => cache !== null) ?? null;
      } else {
        return this.cacheStorage
          .restoreFileStream(pngKey)
          .then((fileStream) => (fileStream ? { fileStream, contentType: 'image/png' as SupportedImageType } : null));
      }
    });
    if (cached) {
      return cached;
    }

    const contributors = await runWithTracing('getAllContributors', () =>
      this.contributorsRepository.getAllContributors(repository),
    );

    const pngImage = await runWithTracing('renderImage', () => this.renderer.render(contributors));
    const webpImage = await runWithTracing('imageToWebp', () => sharp(pngImage).webp({}).toBuffer());

    runWithTracing('saveCache', async () => {
      const pngKey = createCacheKey(repository, 'png');
      const webpKey = createCacheKey(repository, 'webp');
      await Promise.all([this.cacheStorage.saveFile(pngKey, pngImage), this.cacheStorage.saveFile(webpKey, webpImage)]);
    });

    if (options.acceptWebp) {
      return {
        fileStream: Readable.from(webpImage),
        contentType: 'image/webp',
      };
    } else {
      return {
        fileStream: Readable.from(pngImage),
        contentType: 'image/png',
      };
    }
  }
}
