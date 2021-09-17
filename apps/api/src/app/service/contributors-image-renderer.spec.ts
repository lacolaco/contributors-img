import { Contributor } from '@lib/core';
import { container } from 'tsyringe';
import { ContributorsImageRenderer } from './contributors-image-renderer';

describe('ContributorsImageRenderer', () => {
  const service = container.resolve(ContributorsImageRenderer);

  describe('.renderSvg()', () => {
    it('should return SVG string', async () => {
      expect(typeof service.render([]) === 'string');
    });

    it('should render contributors avatar', async () => {
      const contributors: Contributor[] = [
        {
          id: 1,
          avatar_url: 'https://github.com/lacolaco.png',
          login: 'lacolaco',
          html_url: 'https://github.com/lacolaco',
          contributions: 0,
        },
      ];

      const svg = await service.render(contributors);
      expect(svg).toBeTruthy();
    });
  });
});
