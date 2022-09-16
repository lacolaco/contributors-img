import { Injectable } from '@angular/core';
import { AngularFirestore } from '@angular/fire/compat/firestore';
import { filter, map, Observable } from 'rxjs';
import { environment } from '../../../../environments/environment';
import { FeaturedRepository } from '../../../shared/model/repository';
import { FeaturedRepositoryDatastore } from './featured-repository-datastore';

type FeaturedRepositoryDocument = {
  items: FeaturedRepository[];
};

@Injectable()
export class FeaturedRepositoryDatastoreImpl implements FeaturedRepositoryDatastore {
  readonly repositories$: Observable<FeaturedRepository[]> = this.firestore
    .collection(`${environment.firestoreRootCollectionName}`)
    .doc<FeaturedRepositoryDocument>('featured_repositories')
    .valueChanges()
    .pipe(
      filter((doc): doc is FeaturedRepositoryDocument => doc != null),
      map((doc) => doc.items),
    );

  constructor(private readonly firestore: AngularFirestore) {}
}
