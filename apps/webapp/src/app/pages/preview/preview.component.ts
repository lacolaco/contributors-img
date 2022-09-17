import { CommonModule } from '@angular/common';
import { Component, OnDestroy, OnInit } from '@angular/core';
import { ActivatedRoute, Router } from '@angular/router';
import { finalize, map, Subject, switchMap, takeUntil } from 'rxjs';
import { Repository } from '../../models/repository';
import { ContributorsImageApi } from '../../shared/api/contributors-image';
import { FooterComponent } from './footer/footer.component';
import { HeaderComponent } from './header/header.component';
import { ImagePreviewComponent } from './image-preview/image-preview.component';
import { PreviewStore } from './preview.store';
import { RecentUsageComponent } from './recent-usage/recent-usage.component';

@Component({
  templateUrl: './preview.component.html',
  styleUrls: ['./preview.component.scss'],
  standalone: true,
  imports: [CommonModule, HeaderComponent, FooterComponent, ImagePreviewComponent, RecentUsageComponent],
  providers: [PreviewStore],
})
export class PreviewPageComponent implements OnInit, OnDestroy {
  constructor(
    private readonly route: ActivatedRoute,
    private readonly router: Router,
    private readonly api: ContributorsImageApi,
    private readonly store: PreviewStore,
  ) {}

  private readonly onDestroy$ = new Subject<void>();

  ngOnInit() {
    this.route.queryParamMap
      .pipe(
        takeUntil(this.onDestroy$),
        map((q) => ({
          repository: Repository.fromString(q.get('repo') ?? 'angular/angular-ja'),
          max: Number(q.get('max')) || null,
          columns: Number(q.get('columns')) || null,
        })),
        switchMap(({ repository, max, columns }) => {
          this.store.startFetchingImage(repository);
          this.store.closeImageSnippet();
          return this.api
            .getByRepository(repository, { max, columns })
            .pipe(finalize(() => this.store.finishFetchingImage()));
        }),
      )
      .subscribe((imageSvg) => {
        this.store.setImageData(imageSvg);
      });

    this.store
      .select((s) => s.repository)
      .pipe(takeUntil(this.onDestroy$))
      .subscribe((repository) => {
        this.router.navigate([], { queryParams: { repo: repository?.toString() } });
      });
  }

  ngOnDestroy() {
    this.onDestroy$.next();
  }

  selectRepository(repoName: string) {
    this.router.navigate([], { queryParams: { repo: repoName } });
  }
}
