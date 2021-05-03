import { Component, OnDestroy, OnInit } from '@angular/core';
import { ActivatedRoute, Router } from '@angular/router';
import { Repository } from '@lib/core';
import { BehaviorSubject, combineLatest, Subject } from 'rxjs';
import { finalize, map, switchMap, takeUntil } from 'rxjs/operators';
import { ContributorsImageApi } from '../../shared/api/contributors-image';
import { PreviewStore } from './preview.store';

@Component({
  selector: 'app-preview',
  templateUrl: './preview.component.html',
  styleUrls: ['./preview.component.scss'],
})
export class PreviewComponent implements OnInit, OnDestroy {
  constructor(
    private route: ActivatedRoute,
    private router: Router,
    private readonly api: ContributorsImageApi,
    private store: PreviewStore,
  ) {}

  private readonly showImageSnippetSubject = new BehaviorSubject<boolean>(false);
  private readonly onDestroy$ = new Subject();

  readonly state$ = combineLatest([this.store.valueChanges, this.showImageSnippetSubject.asObservable()]).pipe(
    map(([state, showImageSnippet]) => ({
      repository: state.repository,
      imageSvg: state.image.data,
      loading: state.image.fetching > 0,
      showImageSnippet,
    })),
  );

  ngOnInit() {
    this.route.queryParamMap
      .pipe(
        takeUntil(this.onDestroy$),
        map((q) => Repository.fromString(q.get('repo') ?? 'angular/angular-ja')),
        switchMap((repository) => {
          this.store.startFetchingImage(repository);
          return this.api.getByRepository(repository).pipe(finalize(() => this.store.finishFetchingImage()));
        }),
      )
      .subscribe((imageSvg) => {
        this.store.setImageData(imageSvg);
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
