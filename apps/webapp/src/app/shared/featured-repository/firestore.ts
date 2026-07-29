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

    // The listener is registered outside the Angular zone: zone.js patches the transport,
    // so a listener opened inside the zone would make every long-poll keep-alive and retry
    // timer schedule a change-detection pass for the life of the page. Only actual snapshot
    // deliveries re-enter the zone, via the `run` calls below.
    this.repositories$ = new Observable<FeaturedRepositoryDocument | undefined>((subscriber) =>
      ngZone.runOutsideAngular(() =>
        onSnapshot(
          docRef,
          (snapshot) => ngZone.run(() => subscriber.next(snapshot.data())),
          (error) => ngZone.run(() => subscriber.error(error)),
        ),
      ),
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
