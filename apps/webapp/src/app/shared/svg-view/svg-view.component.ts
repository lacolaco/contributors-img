import { CommonModule } from '@angular/common';
import { ChangeDetectionStrategy, Component, Input, NgModule, ViewEncapsulation } from '@angular/core';
import { DomSanitizer } from '@angular/platform-browser';

@Component({
  selector: 'app-svg-view',
  template: `<div *ngIf="sanitizedSvg" [innerHtml]="sanitizedSvg"></div>`,
  styles: [
    `
      app-svg-view svg {
        max-width: 100%;
        height: auto;
      }
    `,
  ],
  changeDetection: ChangeDetectionStrategy.OnPush,
  encapsulation: ViewEncapsulation.None,
})
export class SvgViewComponent {
  @Input() content: string | null = null;

  get sanitizedSvg() {
    if (this.content === null) {
      return null;
    }
    return this.domSanitizer.bypassSecurityTrustHtml(this.content);
  }

  constructor(private readonly domSanitizer: DomSanitizer) {}
}

@NgModule({
  imports: [CommonModule],
  declarations: [SvgViewComponent],
  exports: [SvgViewComponent],
})
export class SvgViewModule {}
