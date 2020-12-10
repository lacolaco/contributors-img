import { Repository } from '@lib/core';
import { Bucket, CacheStorage } from './utils/cache-storage';

export class ContributorsImageCache {
  private readonly storage: CacheStorage;
  constructor(bucket: Bucket, readonly config: { useCache: boolean }) {
    this.storage = new CacheStorage(bucket, config);
  }

  generateCacheFileName(repository: Repository) {
    return `image-cache/${repository.owner}--${repository.repo}`;
  }

  async restore(repository: Repository): Promise<Buffer | null> {
    return this.storage.restore(this.generateCacheFileName(repository), {});
  }

  async save(repository: Repository, file: Buffer): Promise<void> {
    return this.storage.save(this.generateCacheFileName(repository), file, {
      contentType: 'image/png',
    });
  }
}
