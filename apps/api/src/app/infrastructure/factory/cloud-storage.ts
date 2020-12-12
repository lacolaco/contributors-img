import { Storage, Bucket } from '@google-cloud/storage';
import { environment } from '../../../environments/environment';

export const createBucket = (): Bucket | null => {
  if (environment.cacheStorageBucketName) {
    const storage = new Storage();
    return storage.bucket(environment.cacheStorageBucketName);
  } else {
    return null;
  }
};
