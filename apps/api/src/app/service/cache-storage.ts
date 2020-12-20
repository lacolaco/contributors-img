import { Bucket } from '@google-cloud/storage';
import { Readable } from 'stream';
import { inject, injectable } from 'tsyringe';

@injectable()
export class CacheStorage {
  constructor(@inject(Bucket) private readonly bucket: Bucket | null) {}

  async restoreFileStream(filename: string): Promise<Readable | null> {
    if (this.bucket == null) {
      return null;
    }
    const file = this.bucket.file(filename);
    return file.exists().then(([exists]) => {
      if (!exists) {
        return null;
      }
      return file.createReadStream();
    });
  }

  async restoreJSON<T>(filename: string): Promise<T | null> {
    if (this.bucket == null) {
      return null;
    }
    const file = this.bucket.file(filename);
    return file
      .download()
      .then(([data]) => JSON.parse(data.toString()))
      .catch(() => null);
  }

  async saveFile(filename: string, file: Buffer | string): Promise<void> {
    if (this.bucket == null) {
      return;
    }
    await this.bucket.file(filename).save(file, {});
  }

  async saveJSON(filename: string, json: unknown): Promise<void> {
    if (this.bucket == null) {
      return;
    }
    await this.bucket.file(filename).save(JSON.stringify(json), {
      contentType: 'application/json',
    });
  }
}
