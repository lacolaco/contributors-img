import { Contributor } from '@lib/core';
import { createMockContributor } from '@lib/core/testing';
import { createRenderer } from './renderer';

describe('renderer-rust', () => {
  test('returned value is a SVG string', async () => {
    const renderer = createRenderer({ maxColumns: 12 });
    const contributors: Contributor[] = [
      createMockContributor({
        id: 1,
        login: 'login1',
        avatar_url: 'https://via.placeholder.com',
      }),
      createMockContributor({
        id: 2,
        login: 'login2',
        avatar_url: 'https://via.placeholder.com',
      }),
    ];
    const image = await renderer.render(contributors);

    expect(image).toMatchSnapshot();
  });
});
