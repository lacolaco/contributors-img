export type Environment = {
  production: boolean;
  webappUrl: string;
  useHeadless: boolean;
  cacheStorageBucketName: string | null;
  firestoreRootCollectionName: string | null;
};
