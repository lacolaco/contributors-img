import { Firestore } from '@google-cloud/firestore';
import { Logging } from '@google-cloud/logging';
import { Bucket } from '@google-cloud/storage';
import { Octokit } from '@octokit/rest';
import { DependencyContainer } from 'tsyringe';
import { createLogging } from './factory/cloud-logging';
import { createBucket } from './factory/cloud-storage';
import { createFirestore } from './factory/firestore';
import { createOctokit } from './factory/octokit';

export function provideInfrastructure(injector: DependencyContainer) {
  injector.register(Bucket, { useFactory: createBucket });
  injector.register(Octokit, { useFactory: createOctokit });
  injector.register(Firestore, { useFactory: createFirestore });
  injector.register(Logging, { useFactory: createLogging });
}
