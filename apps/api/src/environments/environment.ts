import { Environment } from './type';

export const environment: Environment = {
  production: false,
  cacheStorageBucketName: null,
  firestoreRootCollectionName: 'staging',
  githubAuthToken: process.env.GITHUB_AUTH_TOKEN ?? null,
};
