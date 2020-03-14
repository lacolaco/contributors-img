import { Repository } from '../model/repository';
import * as cacheStorage from './utils/cache-storage';

function generateCacheId(repository: Repository) {
  return `image-cache--${repository.owner}--${repository.repo}`;
}

export async function restoreImageCache(repository: Repository): Promise<Buffer | null> {
  return cacheStorage.readFile(generateCacheId(repository));
}

export async function saveImageCache(repository: Repository, file: Buffer): Promise<void> {
  await cacheStorage.writeFile(generateCacheId(repository), file);
}
