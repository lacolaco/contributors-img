import { Repository, RepoInfo } from '@contributors-img/api-interfaces';

export class RepoInfoRepository {
  constructor(private firestore: FirebaseFirestore.Firestore) {}

  set(repository: Repository, info: Omit<RepoInfo, 'name'>) {
    const repoInfo = this.firestore.collection('repositories').doc(`${repository.owner}--${repository.repo}`);
    return repoInfo.set({
      name: repository.toString(),
      ...info,
    });
  }
}
