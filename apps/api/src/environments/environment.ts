import { Environment } from './type';

export const environment: Environment = {
  production: false,
  cacheStorageBucketName: null,
  firestoreRootCollectionName: 'staging',
  environmentName: 'development',
  githubAuthToken: process.env.GITHUB_AUTH_TOKEN ?? null,
};
