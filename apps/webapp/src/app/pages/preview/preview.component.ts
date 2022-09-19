import { CommonModule } from '@angular/common';
import { Component, inject, OnDestroy, OnInit } from '@angular/core';
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

type PreviewPageParams = {
  repo?: string;
  max?: string;
  columns?: string;
};

@Component({
  templateUrl: './preview.component.html',
  styleUrls: ['./preview.component.scss'],
  standalone: true,
  imports: [CommonModule, HeaderComponent, FooterComponent, ImagePreviewComponent, RecentUsageComponent],
  providers: [PreviewState],
})
export class PreviewPageComponent implements OnInit, OnDestroy {
  private readonly onDestroy$ = new Subject<void>();
  private readonly state = inject(PreviewState);
  private readonly router = inject(Router);
  private readonly imageApi = inject(ContributorsImageApi);
  private readonly queryParams$ = inject(ActivatedRoute).queryParams.pipe(takeUntil(this.onDestroy$));

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
      this.state.patchImageParams({
        repository: repo ? Repository.fromString(repo) : defaultImageParams.repository,
        max: Number(max) || defaultImageParams.max,
        columns: Number(columns) || defaultImageParams.columns,
      });
    });
  }

  private refreshOnStateChange() {
    this.state
      .select('imageParams')
      .pipe(takeUntil(this.onDestroy$))
      .subscribe(async (imageParams) => {
        window.scrollTo({ top: 0, behavior: 'smooth' });
        this.updateQueryParams(imageParams);
        try {
          this.state.startFetchingImage();
          const image = await firstValueFrom(this.imageApi.getImage(imageParams).pipe(delay(100)));
          this.state.finishFetchingImage({ data: image });
        } catch {
          this.state.finishFetchingImage(null);
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
