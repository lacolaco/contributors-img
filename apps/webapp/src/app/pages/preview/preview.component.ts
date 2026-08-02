import { Component, computed, inject, OnDestroy, OnInit, ChangeDetectionStrategy } from '@angular/core';
import { toObservable } from '@angular/core/rxjs-interop';
import { ActivatedRoute, Router } from '@angular/router';
import { delay, firstValueFrom, Subject, takeUntil } from 'rxjs';
import { ImageParams } from '../../models/image-params';
import { Repository } from '../../models/repository';
import { ContributorsImageApi } from '../../shared/api/contributors-image';
import { FooterComponent } from './footer/footer.component';
import { HeaderComponent } from './header/header.component';
import { ImagePreviewComponent } from './image-preview/image-preview.component';
import { RecentUsageComponent } from './recent-usage/recent-usage.component';
import { defaultImageParams, PreviewState } from './state';

interface PreviewPageParams {
  repo?: string;
  max?: string;
  columns?: string;
}

@Component({
  template: `
    <main>
      <app-header />
      <app-image-preview />
      <app-recent-usage />
      <app-footer />
    </main>
  `,
  styleUrls: ['./preview.component.scss'],
  imports: [HeaderComponent, FooterComponent, ImagePreviewComponent, RecentUsageComponent],
  // Angular 22 made OnPush the default and its migration wrote this out to preserve the previous behaviour.
  // Switching to OnPush is a change in behaviour, not part of the upgrade, so it is deferred.
  // eslint-disable-next-line @angular-eslint/prefer-on-push-component-change-detection
  changeDetection: ChangeDetectionStrategy.Eager,
  providers: [PreviewState],
})
export class PreviewPageComponent implements OnInit, OnDestroy {
  private readonly onDestroy$ = new Subject<void>();
  private readonly store = inject(PreviewState);
  private readonly router = inject(Router);
  private readonly imageApi = inject(ContributorsImageApi);
  private readonly queryParams$ = inject(ActivatedRoute).queryParams.pipe(takeUntil(this.onDestroy$));
  private readonly imageParams$ = toObservable(computed(() => this.store.state().imageParams));

  ngOnInit() {
    this.updateStateOnQueryParamsChange();
    this.refreshOnStateChange();
  }

  ngOnDestroy(): void {
    this.onDestroy$.next();
  }

  private updateStateOnQueryParamsChange() {
    this.queryParams$.subscribe((params) => {
      const { repo = null, max = null, columns = null } = params;
      this.store.patchImageParams({
        repository: repo ? Repository.fromString(repo) : defaultImageParams.repository,
        max: Number(max) || defaultImageParams.max,
        columns: Number(columns) || defaultImageParams.columns,
      });
    });
  }

  private refreshOnStateChange() {
    this.imageParams$.pipe(takeUntil(this.onDestroy$)).subscribe(async (imageParams) => {
      window.scrollTo({ top: 0, behavior: 'smooth' });
      this.updateQueryParams(imageParams);
      try {
        this.store.startFetchingImage();
        const image = await firstValueFrom(this.imageApi.getImage(imageParams).pipe(delay(100)));
        this.store.finishFetchingImage({ data: image });
      } catch {
        this.store.finishFetchingImage(null);
      }
    });
  }

  private async updateQueryParams(params: ImageParams) {
    const pageParams: PreviewPageParams = {
      repo: params.repository.toString(),
      max: params.max?.toString() ?? undefined,
      columns: params.columns?.toString() ?? undefined,
    };
    await this.router.navigate([], { queryParams: pageParams });
  }
}
