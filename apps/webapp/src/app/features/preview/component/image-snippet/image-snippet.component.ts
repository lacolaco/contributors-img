import { ChangeDetectionStrategy, Component, Input } from '@angular/core';
import { Repository } from '@lib/core';

@Component({
  selector: 'app-image-snippet',
  templateUrl: './image-snippet.component.html',
  styleUrls: ['./image-snippet.component.scss'],
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class ImageSnippetComponent {
  @Input() repository: Repository;

  get imageSnippet(): string {
    const repoString = this.repository.toString();
    return `
<a href="https://github.com/${repoString}/graphs/contributors">
  <img src="${location.origin}/image?repo=${repoString}" />
</a>

Made with [contrib.rocks](${location.origin}).`.trim();
  }
}
