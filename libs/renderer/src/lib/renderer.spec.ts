import { Contributor } from '@lib/core';
import { createMockContributor } from '@lib/core/testing';
import { Circle } from '@svgdotjs/svg.js';
import { ContributorsImageRenderer, createContributorView, createJsRenderer, createRustRenderer } from './renderer';
import { setupSvgRenderer } from './utils';

describe('renderer', () => {
  let renderer: ContributorsImageRenderer;
  beforeEach(() => {
    setupSvgRenderer();

    renderer = createJsRenderer({ maxColumns: 12 });
  });

  describe('render()', () => {
    test('returned value is a SVG string', async () => {
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

      expect(image).toMatchInlineSnapshot(
        `"<svg xmlns=\\"http://www.w3.org/2000/svg\\" version=\\"1.1\\" xmlns:xlink=\\"http://www.w3.org/1999/xlink\\" xmlns:svgjs=\\"http://svgjs.dev/svgjs\\" width=\\"132\\" height=\\"64\\"><svg width=\\"64\\" height=\\"64\\" x=\\"0\\" y=\\"0\\"><title>login1</title><circle r=\\"31.5\\" cx=\\"32\\" cy=\\"32\\" stroke-width=\\"1\\" stroke=\\"#c0c0c0\\" fill=\\"url(&quot;#SvgjsPattern1000&quot;)\\"></circle><defs><pattern x=\\"0\\" y=\\"0\\" width=\\"64\\" height=\\"64\\" patternUnits=\\"userSpaceOnUse\\" id=\\"SvgjsPattern1000\\"><image width=\\"64\\" height=\\"64\\" href=\\"https://via.placeholder.com\\"></image></pattern></defs></svg><svg width=\\"64\\" height=\\"64\\" x=\\"68\\" y=\\"0\\"><title>login2</title><circle r=\\"31.5\\" cx=\\"32\\" cy=\\"32\\" stroke-width=\\"1\\" stroke=\\"#c0c0c0\\" fill=\\"url(&quot;#SvgjsPattern1001&quot;)\\"></circle><defs><pattern x=\\"0\\" y=\\"0\\" width=\\"64\\" height=\\"64\\" patternUnits=\\"userSpaceOnUse\\" id=\\"SvgjsPattern1001\\"><image width=\\"64\\" height=\\"64\\" href=\\"https://via.placeholder.com\\"></image></pattern></defs></svg></svg>"`,
      );
    });

    test('all given contributors are rendered', async () => {
      const contributors: Contributor[] = [
        { id: 1, login: 'login1', avatar_url: 'https://via.placeholder.com', html_url: 'htmlUrl1', contributions: 1 },
        { id: 2, login: 'login2', avatar_url: 'https://via.placeholder.com', html_url: 'htmlUrl2', contributions: 1 },
      ];

      const image = await renderer.render(contributors);

      expect(contributors.map((c) => c.login).every((login) => image.includes(login))).toBeTruthy();
    });
  });

  describe('createContributorView', () => {
    test('the view size is 64px by default', async () => {
      const svg = await createContributorView(createMockContributor(), { size: 64 });
      expect(svg.width()).toBe(64);
      expect(svg.height()).toBe(64);
    });

    test('the view size is configurable with options', async () => {
      const svg = await createContributorView(createMockContributor(), { size: 48 });
      expect(svg.width()).toBe(48);
      expect(svg.height()).toBe(48);
    });

    test('the view contains a title with contributor name', async () => {
      const contributor = createMockContributor();
      const svg = await createContributorView(contributor, { size: 64 });
      const [title] = svg.children();
      expect(title.svg()).toBe(`<title>${contributor.login}</title>`);
    });

    test('the view shape is a circle', async () => {
      const svg = await createContributorView(createMockContributor(), { size: 64 });
      const [_, shape] = svg.children();
      expect(shape).toBeInstanceOf(Circle);
    });

    test('the shape is filled by image', async () => {
      const contributor = createMockContributor();
      const svg = await createContributorView(contributor, { size: 64 });
      const [_, shape] = svg.children();
      const fill = svg.findOne('pattern');
      expect(shape.fill).toBeDefined();
      // eslint-disable-next-line @typescript-eslint/no-non-null-assertion
      expect(shape.attr('fill')).toContain(fill!.attr('id'));
    });

    test('the image source is contributor avatar', async () => {
      const contributor = createMockContributor();
      const svg = await createContributorView(contributor, { size: 64 });

      // eslint-disable-next-line @typescript-eslint/no-non-null-assertion
      expect(svg.findOne('pattern > image')!.attr('href')).toBe(contributor.avatar_url);
    });
  });
});

describe('renderer-rust', () => {
  test('returned value is a SVG string', async () => {
    const renderer = createRustRenderer({ maxColumns: 12 });
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
