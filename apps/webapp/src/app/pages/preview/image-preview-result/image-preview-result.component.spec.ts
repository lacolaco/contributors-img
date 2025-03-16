import { provideHttpClient } from '@angular/common/http';
import { provideHttpClientTesting } from '@angular/common/http/testing';
import { render, screen } from '@testing-library/angular';
import { Repository } from '../../../models';
import { ImagePreviewResultComponent } from './image-preview-result.component';

describe('ImagePreviewResultComponent', () => {
  it('should render', async () => {
    await render(ImagePreviewResultComponent, {
      providers: [provideHttpClient(), provideHttpClientTesting()],
      componentInputs: {
        repository: Repository.fromString('foo/bar'),
        imageSvg: 'foo',
      },
    });

    expect(screen.getByText('foo')).toBeTruthy();
  });
});
