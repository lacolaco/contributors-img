import { ClipboardModule } from '@angular/cdk/clipboard';
import { MatSnackBarModule } from '@angular/material/snack-bar';
import { render } from '@testing-library/angular';
import { Repository } from '../../../models/repository';
import { ImageSnippetComponent } from './image-snippet.component';

describe('ImageSnippetComponent', () => {
  it('should render', async () => {
    const { fixture } = await render(ImageSnippetComponent, {
      imports: [ClipboardModule, MatSnackBarModule],
      componentInputs: {
        repository: Repository.fromString('foo/bar'),
      },
    });

    const repoString = 'foo/bar';
    const expectedSnippet = `
<a href="https://github.com/${repoString}/graphs/contributors">
  <img src="${location.origin}/image?repo=${repoString}" />
</a>

Made with [contrib.rocks](${location.origin}).`.trim();
    expect(fixture.componentInstance.imageSnippet).toEqual(expectedSnippet);
  });
});
