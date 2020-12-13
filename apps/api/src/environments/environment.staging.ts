import { Environment } from './type';

export const environment: Environment = {
  production: true,
  webappUrl: 'https://stg.contrib.rocks',
  useHeadless: true,
  cacheStorageBucketName: 'staging.contributors-img.appspot.com',
  firestoreRootCollectionName: 'staging',
};
