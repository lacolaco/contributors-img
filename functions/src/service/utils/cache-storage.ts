import { Bucket } from '@google-cloud/storage';
import * as firebase from 'firebase-admin';

async function useBucket<T>(callback: (b: Bucket) => Promise<T>): Promise<T | null> {
  try {
    const bucket = firebase.storage().bucket();
    return await callback(bucket);
  } catch {
    return null;
  }
}

export async function readFile(fileName: string): Promise<Buffer | null> {
  return await useBucket(async bucket => {
    const file = bucket.file(fileName);

    return await file.exists().then(([exists]) => {
      if (exists) {
        return file.download({}).then(([data]) => data);
      }
      return null;
    });
  });
}

export async function writeFile(filename: string, file: Buffer): Promise<void> {
  await useBucket(async bucket => {
    await bucket.file(filename).save(file, {
      public: true,
    });
  });
}
