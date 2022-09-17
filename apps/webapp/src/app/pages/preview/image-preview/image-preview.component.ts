import { CommonModule } from '@angular/common';
import { ChangeDetectionStrategy, Component, inject } from '@angular/core';
import { ImagePreviewFormComponent } from '../image-preview-form/image-preview-form.component';
import { ImagePreviewResultComponent } from '../image-preview-result/image-preview-result.component';
import { PreviewStore } from '../preview.store';

@Component({
  selector: 'app-image-preview',
  template: `
    <ng-container *ngIf="state$ | async as state">
      <app-image-preview-form [repository]="state.repository"> </app-image-preview-form>

      <ng-container *ngIf="state.loading; else showResult">
        <img height="100" src="assets/images/loading.gif" />
      </ng-container>
      <ng-template #showResult>
        <ng-container *ngIf="state.repository && state.imageSvg; else showNoResult">
          <app-image-preview-result [repository]="state.repository" [imageSvg]="state.imageSvg">
          </app-image-preview-result>
        </ng-container>
      </ng-template>
      <ng-template #showNoResult>
        <div>No Result. Is the repository name correct?</div>
      </ng-template>
    </ng-container>
  `,
  styleUrls: ['./image-preview.component.scss'],
  changeDetection: ChangeDetectionStrategy.OnPush,
  standalone: true,
  imports: [CommonModule, ImagePreviewFormComponent, ImagePreviewResultComponent],
})
export class ImagePreviewComponent {
  private readonly store = inject(PreviewStore);

  readonly state$ = this.store.select((state) => ({
    repository: state.repository,
    imageSvg: state.image.data,
    loading: state.image.fetching > 0,
    showImageSnippet: state.showImageSnippet,
  }));
}
