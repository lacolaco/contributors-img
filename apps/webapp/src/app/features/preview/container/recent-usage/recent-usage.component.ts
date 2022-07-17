import { Component, OnDestroy } from '@angular/core';
import { FeaturedRepository } from '@lib/core';
import { Observable, Subject, takeUntil } from 'rxjs';
import { FeaturedRepositoryDatastore } from '../../service/featured-repository-datastore';

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

  readonly repositories$: Observable<FeaturedRepository[]> = this.featuredRepositories.repositories$.pipe(
    takeUntil(this.onDestroy$),
  );

  constructor(private readonly featuredRepositories: FeaturedRepositoryDatastore) {}

  ngOnDestroy() {
    this.onDestroy$.next();
  }
}
