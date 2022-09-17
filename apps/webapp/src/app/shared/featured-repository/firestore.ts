import { Injectable, Provider } from '@angular/core';
import { doc, docData, DocumentReference, Firestore } from '@angular/fire/firestore';
import { filter, map, Observable } from 'rxjs';
import { environment } from '../../../environments/environment';
import { FeaturedRepository } from '../model';
import { FeaturedRepositoryDatasource, FeaturedRepositoryDatasourceToken } from './index';

type FeaturedRepositoryDocument = {
  items: FeaturedRepository[];
};

@Injectable()
export class FirebaseFeaturedRepositoryDatasource implements FeaturedRepositoryDatasource {
  readonly repositories$: Observable<FeaturedRepository[]>;

  constructor(firestore: Firestore) {
    this.repositories$ = docData(
      doc(
        firestore,
        `${environment.firestoreRootCollectionName}/featured_repositories`,
      ) as DocumentReference<FeaturedRepositoryDocument>,
    ).pipe(
      filter((doc): doc is FeaturedRepositoryDocument => doc != null),
      map((doc) => doc.items),
    );
  }
}

export function provideFeaturedRepositoryDatasource(): Provider[] {
  return [
    {
      provide: FeaturedRepositoryDatasourceToken,
      useClass: FirebaseFeaturedRepositoryDatasource,
    },
  ];
}
