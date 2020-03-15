import { Component, OnDestroy, OnInit } from '@angular/core';
import { Observable, Subject } from 'rxjs';
import { takeUntil } from 'rxjs/operators';
import { FetchContributorsUsecase } from './usecase/fetch-contributors.usecase';
import { GitHubContributor } from './core/models';
import { ContributorsStore } from './state/contributors';

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.scss'],
})
export class AppComponent implements OnInit, OnDestroy {
  contributors$: Observable<GitHubContributor[]>;
  contributorsLoading$: Observable<boolean>;
  imageSnippet$: Observable<string>;
  private onDestroy$ = new Subject();

  constructor(private contributorsService: FetchContributorsUsecase, private contributorsStore: ContributorsStore) {
    this.contributors$ = this.contributorsStore.select(state => state.contributors.items);
    this.contributorsLoading$ = this.contributorsStore.select(state => state.contributors.fetching > 0);
    this.imageSnippet$ = this.contributorsStore.select(state => state.imageSnippet);
  }

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
