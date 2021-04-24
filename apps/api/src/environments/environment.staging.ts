import { Environment } from './type';

export const environment: Environment = {
  production: true,
  cacheStorageBucketName: process.env.CACHE_STORAGE_BUCKET ?? null,
  firestoreRootCollectionName: 'staging',
  githubAuthToken: process.env.GITHUB_AUTH_TOKEN ?? null,
};
