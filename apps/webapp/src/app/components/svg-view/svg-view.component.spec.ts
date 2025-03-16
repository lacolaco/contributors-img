import { render, screen } from '@testing-library/angular';
import { SvgViewComponent } from './svg-view.component';

describe('SvgViewComponent', () => {
  it('should render', async () => {
    await render(SvgViewComponent, {
      componentInputs: {
        content: 'test svg',
      },
    });
    expect(screen.getByText('test svg')).toBeTruthy();
  });
});
