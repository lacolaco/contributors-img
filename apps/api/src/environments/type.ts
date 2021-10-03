export type Environment = {
  production: boolean;
  cacheStorageBucketName: string | null;
  environmentName: string;
  firestoreRootCollectionName: string | null;
  githubAuthToken: string | null;
};
