import { Component, OnDestroy, OnInit } from '@angular/core';
import { ActivatedRoute, Router } from '@angular/router';
import { Repository } from '@lib/core';
import { Subject } from 'rxjs';
import { finalize, map, switchMap, takeUntil } from 'rxjs/operators';
import { ContributorsImageApi } from '../../shared/api/contributors-image';
import { PreviewStore } from './preview.store';

@Component({
  selector: 'app-preview',
  templateUrl: './preview.component.html',
  styleUrls: ['./preview.component.scss'],
  providers: [PreviewStore],
})
export class PreviewComponent implements OnInit, OnDestroy {
  constructor(
    private readonly route: ActivatedRoute,
    private readonly router: Router,
    private readonly api: ContributorsImageApi,
    private readonly store: PreviewStore,
  ) {}

  private readonly onDestroy$ = new Subject<void>();

  readonly state$ = this.store.select((state) => ({
    repository: state.repository,
    imageSvg: state.image.data,
    loading: state.image.fetching > 0,
    showImageSnippet: state.showImageSnippet,
  }));

  ngOnInit() {
    this.route.queryParamMap
      .pipe(
        takeUntil(this.onDestroy$),
        map((q) => ({
          repository: Repository.fromString(q.get('repo') ?? 'angular/angular-ja'),
          max: Number(q.get('max')) || null,
        })),
        switchMap(({ repository, max }) => {
          this.store.startFetchingImage(repository);
          this.store.closeImageSnippet();
          return this.api.getByRepository(repository, { max }).pipe(finalize(() => this.store.finishFetchingImage()));
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
    this.router.navigate([], { queryParams: { repo: repoName } });
  }

  showImageSnippet() {
    this.store.showImageSnippet();
  }
}
