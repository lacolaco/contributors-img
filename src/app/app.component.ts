import { Component, OnDestroy, OnInit } from '@angular/core';
import { Observable, Subject } from 'rxjs';
import { takeUntil } from 'rxjs/operators';
import { ContributorsService } from './contributors.service';
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

  constructor(private contributorsService: ContributorsService, private contributorsStore: ContributorsStore) {
    this.contributors$ = this.contributorsStore.selectValue(state => state.items);
    this.contributorsLoading$ = this.contributorsStore.selectValue(state => state.itemsLoading > 0);
    this.imageSnippet$ = this.contributorsStore.selectValue(state => state.imageSnippet);
  }

  ngOnInit() {
    this.contributorsStore
      .selectValue(state => state.repository)
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
