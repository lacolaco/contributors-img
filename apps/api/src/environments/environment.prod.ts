import { Environment } from './type';

export const environment: Environment = {
  production: true,
  webappUrl: 'https://contrib.rocks',
  useHeadless: true,
  cacheStorageBucketName: 'contributors-img.appspot.com',
  firestoreRootCollectionName: 'production',
};
