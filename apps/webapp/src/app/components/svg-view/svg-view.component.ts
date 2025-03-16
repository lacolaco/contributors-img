import { ChangeDetectionStrategy, Component, ViewEncapsulation, inject, input } from '@angular/core';
import { DomSanitizer } from '@angular/platform-browser';

@Component({
  selector: 'app-svg-view',
  template: `@if (sanitizedSvg) {
    <div [innerHtml]="sanitizedSvg"></div>
  }`,
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
  imports: [],
})
export class SvgViewComponent {
  private readonly domSanitizer = inject(DomSanitizer);

  readonly content = input<string | null>(null);

  get sanitizedSvg() {
    const content = this.content();
    if (content === null) {
      return null;
    }
    return this.domSanitizer.bypassSecurityTrustHtml(content);
  }
}
