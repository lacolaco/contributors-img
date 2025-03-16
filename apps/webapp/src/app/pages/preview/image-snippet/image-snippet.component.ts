import { Clipboard, ClipboardModule } from '@angular/cdk/clipboard';
import { ChangeDetectionStrategy, Component, inject, Input } from '@angular/core';
import { MatSnackBar, MatSnackBarModule } from '@angular/material/snack-bar';
import { Repository } from '../../../models/repository';

@Component({
  selector: 'app-image-snippet',
  template: `
    <div class="heading">Copy & Paste to README.md</div>
    <textarea class="snippet" readonly rows="5" (click)="copyToClipboard()">{{ imageSnippet }}</textarea>
    <div>
      <span class="">Available options (Add to image URL query params)</span>
      <dl style="margin: 4px 0">
        <dt style="font-family: monospace; font-size: 14px">{{ 'max={number}' }}</dt>
        <dd style="font-size: 14px; margin-inline-start: 2rem">Max contributor count (100 by default)</dd>
        <dt style="font-family: monospace; font-size: 14px">{{ 'columns={number}' }}</dt>
        <dd style="font-size: 14px; margin-inline-start: 2rem">Max columns (12 by default)</dd>
        <dt style="font-family: monospace; font-size: 14px">{{ 'anon={0|1}' }}</dt>
        <dd style="font-size: 14px; margin-inline-start: 2rem">Include anonymous users (false by default)</dd>
      </dl>
    </div>
  `,
  styleUrls: ['./image-snippet.component.scss'],
  changeDetection: ChangeDetectionStrategy.OnPush,
  imports: [ClipboardModule, MatSnackBarModule],
})
export class ImageSnippetComponent {
  private readonly clipboard = inject(Clipboard);
  private readonly snackBar = inject(MatSnackBar);

  @Input() repository: Repository;

  get imageSnippet(): string {
    const repoString = this.repository.toString();
    return `
<a href="https://github.com/${repoString}/graphs/contributors">
  <img src="${location.origin}/image?repo=${repoString}" />
</a>

Made with [contrib.rocks](${location.origin}).`.trim();
  }

  copyToClipboard() {
    this.clipboard.copy(this.imageSnippet);
    this.snackBar.open('Copied to clipboard!', undefined, { duration: 2000 });
  }
}
