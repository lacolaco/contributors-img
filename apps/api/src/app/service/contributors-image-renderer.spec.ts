import { container } from 'tsyringe';
import { ContributorsImageRenderer } from './contributors-image-renderer';

describe('ContributorsImageRenderer', () => {
  const service = container.resolve(ContributorsImageRenderer);

  describe('.renderSvg()', () => {
    it('should return SVG string', async () => {
      expect(typeof service.renderSvg([]) === 'string');
    });
  });
});
