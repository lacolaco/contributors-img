import { render, screen } from '@testing-library/angular';
import { HeaderComponent } from './header.component';

describe('HeaderComponent', () => {
  it('should render', async () => {
    await render(HeaderComponent);

    expect(screen.getByText('contrib.rocks')).toBeTruthy();
  });
});
