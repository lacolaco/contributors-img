import { ChangeDetectionStrategy, Component, OnChanges, input } from '@angular/core';
import { MatButtonModule } from '@angular/material/button';
import { SvgViewComponent } from '../../../components/svg-view/svg-view.component';
import { Repository } from '../../../models';
import { ImageSnippetComponent } from '../image-snippet/image-snippet.component';

@Component({
  selector: 'app-image-preview-result',
  template: `
    <div class="pane">
      <app-svg-view [content]="imageSvg()" />
      @if (snippetOpen) {
        <app-image-snippet [repository]="repository()" />
      }
      @if (!snippetOpen) {
        <button mat-stroked-button color="primary" (click)="showImageSnippet()">Get Image URL!</button>
      }
    </div>
  `,
  styleUrls: ['./image-preview-result.component.scss'],
  changeDetection: ChangeDetectionStrategy.OnPush,
  imports: [MatButtonModule, SvgViewComponent, ImageSnippetComponent],
})
export class ImagePreviewResultComponent implements OnChanges {
  readonly repository = input.required<Repository>();
  readonly imageSvg = input.required<string>();

  snippetOpen = false;

  ngOnChanges() {
    this.snippetOpen = false;
  }

  showImageSnippet() {
    this.snippetOpen = true;
  }
}
