import { CommonModule } from '@angular/common';
import { ChangeDetectionStrategy, Component, inject } from '@angular/core';
import { map } from 'rxjs';
import { ImagePreviewFormComponent } from '../image-preview-form/image-preview-form.component';
import { ImagePreviewResultComponent } from '../image-preview-result/image-preview-result.component';
import { PreviewState } from '../state';

@Component({
  selector: 'app-image-preview',
  template: `
    @if (state$ | async; as state) {
      <app-image-preview-form [value]="state.imageParams" />
      @if (state.loading) {
        <img height="100" src="assets/images/loading.gif" />
      } @else {
        @if (state.result) {
          <app-image-preview-result [repository]="state.imageParams.repository" [imageSvg]="state.result.data" />
        } @else {
          <div>No Result. Is the repository name correct?</div>
        }
      }
    }
  `,
  styleUrls: ['./image-preview.component.scss'],
  changeDetection: ChangeDetectionStrategy.OnPush,
  imports: [CommonModule, ImagePreviewFormComponent, ImagePreviewResultComponent],
})
export class ImagePreviewComponent {
  private readonly state = inject(PreviewState);

  readonly state$ = this.state.select().pipe(
    map((state) => ({
      imageParams: state.imageParams,
      loading: state.fetchingCount > 0,
      result: state.result,
    })),
  );
}
