import { Firestore } from '@google-cloud/firestore';
import { Bucket } from '@google-cloud/storage';
import { Octokit } from '@octokit/rest';
import { container } from 'tsyringe';
import { createBucket } from './factory/cloud-storage';
import { createFirestore } from './factory/firestore';
import { createOctokit } from './factory/octokit';

container.register(Bucket, {
  useFactory: createBucket,
});
container.register(Octokit, {
  useFactory: createOctokit,
});
container.register(Firestore, {
  useFactory: createFirestore,
});
