import { renderContributorsImage } from './renderer';
import { Contributor } from '@lib/core';

describe('renderContributorsImage', () => {
  test('returned value is a SVG string', async () => {
    const contributors: Contributor[] = [];

    const image = await renderContributorsImage(contributors);

    expect(typeof image).toBe('string');
  });

  test('all given contributors are rendered', async () => {
    const contributors: Contributor[] = [
      { id: 1, login: 'login1', avatar_url: 'https://via.placeholder.com', html_url: 'htmlUrl1', contributions: 1 },
      { id: 2, login: 'login2', avatar_url: 'https://via.placeholder.com', html_url: 'htmlUrl2', contributions: 1 },
    ];

    const image = await renderContributorsImage(contributors);

    expect(contributors.map((c) => c.login).every((login) => image.includes(login))).toBeTruthy();
  });
});
