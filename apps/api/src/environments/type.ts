export type Environment = {
  production: boolean;
  useHeadless: boolean;
  cacheStorageBucketName: string | null;
  firestoreRootCollectionName: string | null;
  githubAuthToken: string | null;
};
