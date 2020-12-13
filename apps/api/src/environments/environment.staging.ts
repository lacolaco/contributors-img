import { Environment } from './type';

export const environment: Environment = {
  production: true,
  webappUrl: 'https://contributors-img-dev.web.app',
  useHeadless: true,
  cacheStorageBucketName: 'staging.contributors-img.appspot.com',
  firestoreRootCollectionName: 'staging',
};
