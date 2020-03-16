import { Component, OnDestroy, OnInit } from '@angular/core';
import { Subject } from 'rxjs';
import { takeUntil } from 'rxjs/operators';
import { ContributorsStore } from './state/contributors';
import { FetchContributorsUsecase } from './usecase/fetch-contributors.usecase';

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.scss'],
})
export class AppComponent implements OnInit, OnDestroy {
  constructor(private contributorsService: FetchContributorsUsecase, private contributorsStore: ContributorsStore) {}

  readonly state$ = this.contributorsStore.select(state => ({
    contributors: state.contributors.items,
    loading: state.contributors.fetching > 0,
    imageSnippet: state.imageSnippet,
  }));
  private readonly onDestroy$ = new Subject();

  ngOnInit() {
    this.contributorsStore
      .select(state => state.repository)
      .pipe(takeUntil(this.onDestroy$))
      .subscribe(repository => {
        try {
          this.contributorsService.fetchContributors(repository);
        } catch (err) {
          // console.error(err);
        }
      });
  }

  ngOnDestroy() {
    this.onDestroy$.next();
  }
}
