import { Environment } from './type';

export const environment: Environment = {
  production: true,
  useHeadless: true,
  cacheStorageBucketName: 'staging.contributors-img.appspot.com',
  firestoreRootCollectionName: 'staging',
};
