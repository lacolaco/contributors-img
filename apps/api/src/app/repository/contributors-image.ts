import { Contributor, Repository } from '@lib/core';
import { injectable } from 'tsyringe';
import { CacheStorage } from '../service/cache-storage';
import { ContributorsImageRenderer } from '../service/contributors-image-renderer';
import { runWithTracing } from '../utils/tracing';

function createCacheKey(repository: Repository) {
  return `image-cache/${repository.owner}--${repository.repo}`;
}

@injectable()
export class ContributorsImageRepository {
  constructor(private readonly cacheStorage: CacheStorage, private readonly renderer: ContributorsImageRenderer) {}

  async getImage(repository: Repository, contributors: Contributor[]): Promise<Buffer> {
    const cacheKey = createCacheKey(repository);
    const cached = await runWithTracing('restoreFromCache', () => this.cacheStorage.restore(cacheKey));
    if (cached) {
      return cached;
    }
    const image = await runWithTracing('renderImage', () => this.renderer.render(contributors));

    await runWithTracing('saveCache', () => this.cacheStorage.save(cacheKey, image));
    return image;
  }
}
