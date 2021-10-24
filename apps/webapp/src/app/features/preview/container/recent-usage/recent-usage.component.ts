import { Component, OnDestroy } from '@angular/core';
import { AngularFirestore } from '@angular/fire/compat/firestore';
import { FeaturedRepository } from '@lib/core';
import { filter, map, Observable, Subject, takeUntil } from 'rxjs';
import { environment } from '../../../../../environments/environment';

type FeaturedRepositoryDocument = {
  items: FeaturedRepository[];
};

@Component({
  selector: 'app-recent-usage',
  templateUrl: './recent-usage.component.html',
  styles: [
    `
      :host {
        display: block;
        width: 100%;
      }
    `,
  ],
})
export class RecentUsageComponent implements OnDestroy {
  private readonly onDestroy$ = new Subject<void>();

  readonly repositories$: Observable<FeaturedRepository[]> = this.firestore
    .collection(`${environment.firestoreRootCollectionName}`)
    .doc<FeaturedRepositoryDocument>('featured_repositories')
    .valueChanges()
    .pipe(
      takeUntil(this.onDestroy$),
      filter((doc): doc is FeaturedRepositoryDocument => doc != null),
      map((doc) => doc.items),
    );

  constructor(private readonly firestore: AngularFirestore) {}

  ngOnDestroy() {
    this.onDestroy$.next();
  }
}
