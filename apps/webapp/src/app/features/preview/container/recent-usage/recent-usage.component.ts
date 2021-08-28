import { Component, OnDestroy } from '@angular/core';
import { AngularFirestore } from '@angular/fire/compat/firestore';
import { Subject } from 'rxjs';
import { takeUntil, throttleTime } from 'rxjs/operators';
import { environment } from '../../../../../environments/environment';

@Component({
  selector: 'app-recent-usage',
  templateUrl: './recent-usage.component.html',
  styles: [
    `
      :host {
        display: block;
      }
    `,
  ],
})
export class RecentUsageComponent implements OnDestroy {
  private readonly onDestroy$ = new Subject<void>();

  readonly repositories$ = this.firestore
    .collection<{ name: string }>(`${environment.firestoreRootCollectionName}/usage/repositories`, (q) =>
      q.limit(12).orderBy('timestamp', 'desc'),
    )
    .valueChanges()
    .pipe(takeUntil(this.onDestroy$), throttleTime(1000 * 10));

  constructor(private readonly firestore: AngularFirestore) {}

  ngOnDestroy() {
    this.onDestroy$.next();
  }
}
