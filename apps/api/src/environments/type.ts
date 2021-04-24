export type Environment = {
  production: boolean;
  cacheStorageBucketName: string | null;
  firestoreRootCollectionName: string | null;
  githubAuthToken: string | null;
};
