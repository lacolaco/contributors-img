import { Environment } from './type';

export const environment: Environment = {
  production: true,
  useHeadless: true,
  cacheStorageBucketName: 'staging.contributors-img.appspot.com',
  firestoreRootCollectionName: 'staging',
  githubAuthToken: process.env.GITHUB_AUTH_TOKEN ?? null,
};
