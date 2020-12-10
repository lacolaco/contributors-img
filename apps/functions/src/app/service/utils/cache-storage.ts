// tslint:disable-next-line: no-implicit-dependencies
export type Bucket = import('@google-cloud/storage').Bucket;
export type SaveOptions = import('@google-cloud/storage').SaveOptions;
export type DownloadOptions = import('@google-cloud/storage').DownloadOptions;

export class CacheStorage {
  constructor(private bucket: Bucket, private readonly config: { useCache: boolean }) {}

  async restore(filename: string, options: DownloadOptions = {}): Promise<Buffer | null> {
    if (!this.config.useCache) {
      return null;
    }
    console.log(`CacheStorage: readFile ${filename}`);
    const file = this.bucket.file(filename);
    return await file.exists().then(([exists]) => {
      if (exists) {
        return file.download({ ...options }).then(([data]) => data);
      }
      return null;
    });
  }

  async save(filename: string, file: Buffer | string, options: SaveOptions = {}): Promise<void> {
    if (!this.config.useCache) {
      return;
    }
    console.log(`CacheStorage: writeFile: ${filename}`);
    await this.bucket.file(filename).save(file, {
      public: true,
      ...options,
    });
  }
}
