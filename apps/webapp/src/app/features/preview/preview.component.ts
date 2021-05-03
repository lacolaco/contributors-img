import { Component, OnDestroy, OnInit } from '@angular/core';
import { ActivatedRoute, Router } from '@angular/router';
import { BehaviorSubject, combineLatest, Subject } from 'rxjs';
import { map, takeUntil } from 'rxjs/operators';
import { PreviewStore } from './preview.store';
import { FetchContributorsUsecase } from './usecase/fetch-contributors.usecase';

@Component({
  selector: 'app-preview',
  templateUrl: './preview.component.html',
  styleUrls: ['./preview.component.scss'],
})
export class PreviewComponent implements OnInit, OnDestroy {
  constructor(
    private route: ActivatedRoute,
    private router: Router,
    private fetchContributors: FetchContributorsUsecase,
    private store: PreviewStore,
  ) {}

  private readonly showImageSnippetSubject = new BehaviorSubject<boolean>(false);
  private readonly onDestroy$ = new Subject();

  readonly state$ = combineLatest([this.store.valueChanges, this.showImageSnippetSubject.asObservable()]).pipe(
    map(([state, showImageSnippet]) => ({
      repository: state.repository,
      contributors: state.contributors.items,
      loading: state.contributors.fetching > 0,
      showImageSnippet,
    })),
  );

  ngOnInit() {
    this.route.queryParamMap
      .pipe(
        takeUntil(this.onDestroy$),
        map((q) => q.get('repo')),
      )
      .subscribe((repo) => {
        this.fetchContributors.execute(repo || 'angular/angular-ja');
      });
  }

  ngOnDestroy() {
    this.onDestroy$.next();
  }

  selectRepository(repoName: string) {
    this.showImageSnippetSubject.next(false);
    this.router.navigate([], { queryParams: { repo: repoName } });
  }

  showImageSnippet() {
    this.showImageSnippetSubject.next(true);
  }
}
