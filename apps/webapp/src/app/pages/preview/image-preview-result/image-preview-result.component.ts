import { CommonModule } from '@angular/common';
import { ChangeDetectionStrategy, Component, Input, OnChanges } from '@angular/core';
import { MatButtonModule } from '@angular/material/button';
import { SvgViewComponent } from '../../../components/svg-view/svg-view.component';
import { Repository } from '../../../models';
import { ImageSnippetComponent } from '../image-snippet/image-snippet.component';

@Component({
  selector: 'app-image-preview-result',
  template: `
    <div class="pane">
      <app-svg-view [content]="imageSvg"></app-svg-view>
      <app-image-snippet *ngIf="snippetOpen" [repository]="repository"> </app-image-snippet>
      <button *ngIf="!snippetOpen" mat-stroked-button color="primary" (click)="showImageSnippet()">
        Get Image URL!
      </button>
    </div>
  `,
  styleUrls: ['./image-preview-result.component.scss'],
  changeDetection: ChangeDetectionStrategy.OnPush,
  standalone: true,
  imports: [CommonModule, MatButtonModule, SvgViewComponent, ImageSnippetComponent],
})
export class ImagePreviewResultComponent implements OnChanges {
  @Input() repository: Repository;
  @Input() imageSvg: string;

  snippetOpen = false;

  ngOnChanges() {
    this.snippetOpen = false;
  }

  showImageSnippet() {
    this.snippetOpen = true;
  }
}
