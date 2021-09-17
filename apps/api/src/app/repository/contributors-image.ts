import { Repository } from '@lib/core';
import { injectable } from 'tsyringe';
import { CacheStorage } from '../service/cache-storage';
import { ContentType } from '../utils/content-type';
import { runWithTracing } from '../utils/tracing';
import { FileStream } from '../utils/types';
import { GetContributorsParams } from './contributors';

function createCacheKey(repository: Repository, params: GetContributorsParams, ext: string) {
  return `image-cache/${repository.owner}--${repository.repo}--${params.maxCount}.${ext}`;
}

type SavedImage =
  | FileStream
  | {
      readonly data: null;
      save: (data: string, contentType: ContentType) => Promise<void>;
    };

@injectable()
export class ContributorsImageRepository {
  constructor(private readonly cacheStorage: CacheStorage) {}

  async loadImage(repository: Repository, params: { maxCount: number }): Promise<SavedImage> {
    return runWithTracing('ContributorsImageRepository.loadImage', async () => {
      const cacheKey = createCacheKey(repository, params, 'svg');
      const fileStream = await this.cacheStorage.restoreFileStream(cacheKey);
      if (fileStream) {
        return fileStream;
      }
      return {
        data: null,
        save: async (data, contentType) => {
          await this.cacheStorage.saveFile(cacheKey, data, contentType);
        },
      };
    });
  }
}
