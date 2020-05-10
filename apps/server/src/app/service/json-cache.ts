import { Contributor, Repository } from '@contributors-img/api-interfaces';
import { Bucket, CacheStorage } from './utils/cache-storage';

export class ContributorsJsonCache {
  private readonly storage: CacheStorage;
  constructor(bucket: Bucket, readonly config: { useCache: boolean }) {
    this.storage = new CacheStorage(bucket, config);
  }

  generateCacheFileName(repository: Repository) {
    return `contributors-json-cache/${repository.owner}--${repository.repo}.json`;
  }

  async restore(repository: Repository): Promise<Contributor[] | null> {
    return this.storage.restore(this.generateCacheFileName(repository), {}).then((file) => {
      if (file === null) {
        return null;
      }
      return JSON.parse(file.toString()) as Contributor[];
    });
  }

  async save(repository: Repository, data: Contributor[]): Promise<void> {
    return this.storage.save(this.generateCacheFileName(repository), JSON.stringify(data), {
      contentType: 'application/json',
    });
  }
}
