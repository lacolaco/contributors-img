import { Contributor, Repository } from '@lib/core';
import { injectable } from 'tsyringe';
import { CacheStorage } from '../service/cache-storage';
import { ContributorsImageRenderer } from '../service/contributors-image-renderer';

function createCacheKey(repository: Repository) {
  return `image-cache/${repository.owner}--${repository.repo}`;
}

@injectable()
export class ContributorsImageRepository {
  constructor(private readonly cacheStorage: CacheStorage, private readonly renderer: ContributorsImageRenderer) {}

  async getImage(repository: Repository, contributors: Contributor[]): Promise<Buffer> {
    const cacheKey = createCacheKey(repository);
    const cached = await this.cacheStorage.restore(cacheKey);
    if (cached) {
      return cached;
    }
    const image = await this.renderer.render(contributors);
    await this.cacheStorage.save(cacheKey, image);
    return image;
  }
}
