import { Bucket, DownloadOptions, SaveOptions } from '@google-cloud/storage';
import { injectable, inject } from 'tsyringe';

@injectable()
export class CacheStorage {
  constructor(@inject(Bucket) private readonly bucket: Bucket | null) {}

  async restore(filename: string, options: DownloadOptions = {}): Promise<Buffer | null> {
    if (this.bucket == null) {
      return null;
    }
    const file = this.bucket.file(filename);
    return await file.exists().then(([exists]) => {
      if (!exists) {
        return null;
      }
      return file.download({ ...options }).then(([data]) => data);
    });
  }

  async restoreJSON<T>(filename: string): Promise<T | null> {
    const file = await this.restore(filename);
    if (!file) {
      return null;
    }
    return JSON.parse(file.toString());
  }

  async save(filename: string, file: Buffer | string, options: SaveOptions = {}): Promise<void> {
    if (this.bucket == null) {
      return;
    }
    await this.bucket.file(filename).save(file, {
      public: true,
      ...options,
    });
  }

  async saveJSON(filename: string, json: unknown): Promise<void> {
    this.save(filename, JSON.stringify(json), { contentType: 'application/json' });
  }
}
