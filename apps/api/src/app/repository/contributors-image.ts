import { Repository } from '@lib/core';
import { Readable } from 'stream';
import { injectable } from 'tsyringe';
import { CacheStorage } from '../service/cache-storage';
import { ContributorsImageRenderer } from '../service/contributors-image-renderer';
import { runWithTracing } from '../utils/tracing';
import { SupportedImageType } from '../utils/types';
import { ContributorsRepository } from './contributors';
import { convertToWebp } from '../utils/image';

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
    const webpImage = await runWithTracing('imageToWebp', () => convertToWebp(pngImage));

    runWithTracing('saveCache', async () => {
      const pngKey = createCacheKey(repository, 'png');
      const webpKey = createCacheKey(repository, 'webp');
      await Promise.all([
        this.cacheStorage.saveFile(pngKey, pngImage, 'image/png'),
        this.cacheStorage.saveFile(webpKey, webpImage, 'image/webp'),
      ]);
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

  async getImageFileStream2(
    repository: Repository,
  ): Promise<{ fileStream: Readable; contentType: SupportedImageType }> {
    const cacheKey = createCacheKey(repository, 'svg');
    const cached = await runWithTracing('restoreFromCache', async () => {
      return this.cacheStorage
        .restoreFileStream(cacheKey)
        .then((fileStream) => (fileStream ? { fileStream, contentType: 'image/svg+xml' as SupportedImageType } : null));
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
