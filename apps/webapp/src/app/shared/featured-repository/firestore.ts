import { Injectable, NgZone, Provider, inject } from '@angular/core';
import { doc, DocumentReference, onSnapshot } from 'firebase/firestore';
import { filter, map, Observable } from 'rxjs';
import { environment } from '../../../environments/environment';
import { FeaturedRepository } from '../../models';
import { FirestoreToken } from '../firebase';
import { FeaturedRepositoryDatasource, FeaturedRepositoryDatasourceToken } from './index';

interface FeaturedRepositoryDocument {
  items: FeaturedRepository[];
}

@Injectable()
export class FirebaseFeaturedRepositoryDatasource implements FeaturedRepositoryDatasource {
  readonly repositories$: Observable<FeaturedRepository[]>;

  constructor() {
    const firestore = inject(FirestoreToken);
    const ngZone = inject(NgZone);

    const docRef = doc(
      firestore,
      `${environment.firestoreRootCollectionName}/featured_repositories`,
    ) as DocumentReference<FeaturedRepositoryDocument>;

    this.repositories$ = new Observable<FeaturedRepositoryDocument | undefined>((subscriber) => {
      const unsubscribe = onSnapshot(
        docRef,
        (snapshot) => ngZone.run(() => subscriber.next(snapshot.data())),
        (error) => ngZone.run(() => subscriber.error(error)),
      );
      return unsubscribe;
    }).pipe(
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
