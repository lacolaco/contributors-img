import { provideHttpClient } from '@angular/common/http';
import { provideHttpClientTesting } from '@angular/common/http/testing';
import { render, screen } from '@testing-library/angular';
import { provideNoopFeaturedRepositoryDatasource } from '../../shared/featured-repository/noop';
import { PreviewPageComponent } from './preview.component';

describe('PreviewComponent', () => {
  it('should render', async () => {
    await render(PreviewPageComponent, {
      providers: [provideHttpClient(), provideHttpClientTesting(), provideNoopFeaturedRepositoryDatasource()],
    });

    expect(screen.getByText('Generate an image of contributors to keep your README.md in sync.')).toBeTruthy();
  });
});
