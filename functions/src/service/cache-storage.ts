type StorageBucket = ReturnType<(ReturnType<typeof import('firebase-admin').storage>)['bucket']>;

export async function readFile(bucket: StorageBucket, filename: string): Promise<Buffer | null> {
  console.log(`readFile: ${filename}`);
  try {
    const file = bucket.file(filename);
    return await file.exists().then(([exists]) => {
      if (exists) {
        return file.download({}).then(([data]) => data);
      }
      return null;
    });
  } catch (error) {
    console.error(error);
    return null;
  }
}

export async function writeFile(bucket: StorageBucket, filename: string, file: Buffer): Promise<void> {
  console.log(`writeFile: ${filename}`);
  try {
    await bucket.file(filename).save(file, {
      public: true,
    });
  } catch (error) {
    console.error(error);
  }
}
