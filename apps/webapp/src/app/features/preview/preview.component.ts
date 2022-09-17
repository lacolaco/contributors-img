import { CommonModule } from '@angular/common';
import { Component, OnDestroy, OnInit } from '@angular/core';
import { MatButtonModule } from '@angular/material/button';
import { ActivatedRoute, Router } from '@angular/router';
import { finalize, map, Subject, switchMap, takeUntil } from 'rxjs';
import { ContributorsImageApi } from '../../shared/api/contributors-image';
import { Repository } from '../../shared/model/repository';
import { SvgViewComponent } from '../../shared/svg-view/svg-view.component';
import { FooterComponent } from './component/footer/footer.component';
import { HeaderComponent } from './component/header/header.component';
import { ImageSnippetComponent } from './component/image-snippet/image-snippet.component';
import { RepositoryFormComponent } from './component/repository-form/repository-form.component';
import { RecentUsageComponent } from './container/recent-usage/recent-usage.component';
import { PreviewStore } from './preview.store';

@Component({
  selector: 'app-preview',
  templateUrl: './preview.component.html',
  styleUrls: ['./preview.component.scss'],
  standalone: true,
  imports: [
    CommonModule,
    MatButtonModule,
    HeaderComponent,
    FooterComponent,
    ImageSnippetComponent,
    RepositoryFormComponent,
    SvgViewComponent,
    RecentUsageComponent,
  ],
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
