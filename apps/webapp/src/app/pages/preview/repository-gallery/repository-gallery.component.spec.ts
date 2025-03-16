import { provideRouter } from '@angular/router';
import { render, screen } from '@testing-library/angular';
import { RepositoryGalleryComponent } from './repository-gallery.component';

describe('RepositoryGalleryComponent', () => {
  it('should render', async () => {
    await render(RepositoryGalleryComponent, {
      providers: [provideRouter([])],
      componentInputs: {
        repositories: [{ repository: 'test/repo', stargazers: 123 }],
      },
    });

    expect(screen.getByText('test/repo')).toBeTruthy();
  });
});
