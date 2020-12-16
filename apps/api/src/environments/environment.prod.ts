import { Environment } from './type';

export const environment: Environment = {
  production: true,
  useHeadless: true,
  cacheStorageBucketName: 'contributors-img.appspot.com',
  firestoreRootCollectionName: 'production',
  githubAuthToken: process.env.GITHUB_AUTH_TOKEN ?? null,
};
