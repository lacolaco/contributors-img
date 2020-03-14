import { Repository } from '../model/repository';

// tslint:disable-next-line: no-implicit-dependencies
type Bucket = import('@google-cloud/storage').Bucket;

function generateCacheId(repository: Repository) {
  return `image-cache--${repository.owner}--${repository.repo}`;
}

export class ContributorsImageCache {
  constructor(private bucket: Bucket) {}

  async restore(repository: Repository): Promise<Buffer | null> {
    const filename = generateCacheId(repository);
    console.log(`ContributorsImageCache: readFile ${filename}`);
    const file = this.bucket.file(filename);

    return await file.exists().then(([exists]) => {
      if (exists) {
        return file.download({}).then(([data]) => data);
      }
      return null;
    });
  }

  async save(repository: Repository, file: Buffer): Promise<void> {
    const filename = generateCacheId(repository);
    console.log(`ContributorsImageCache: writeFile: ${filename}`);
    await this.bucket.file(filename).save(file, {
      public: true,
    });
  }
}
