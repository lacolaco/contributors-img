import { Firestore } from '@google-cloud/firestore';
import { Repository } from '@lib/core';
import { injectable } from 'tsyringe';
import { environment } from '../../environments/environment';

@injectable()
export class UsageRepository {
  constructor(private readonly firestore: Firestore) {}

  async saveRepositoryUsage(repository: Repository, timestamp: Date) {
    this.firestore
      .collection(`${environment.firestoreRootCollectionName}/usage/repositories`)
      .doc(`${repository.owner}--${repository.repo}`)
      .set({
        name: repository.toString(),
        timestamp,
      });
  }
}
