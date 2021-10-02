import { Contributor } from '@lib/core';
import { createMockContributor } from '@lib/core/testing';
import { Circle } from '@svgdotjs/svg.js';
import { ContributorsImageRenderer, createContributorView, createRenderer } from './renderer';
import { setupSvgRenderer } from './utils';

describe('renderer', () => {
  let renderer: ContributorsImageRenderer;
  beforeEach(() => {
    setupSvgRenderer();

    renderer = createRenderer();
  });

  describe('render()', () => {
    test('returned value is a SVG string', async () => {
      const contributors: Contributor[] = [];

      const image = await renderer.render(contributors);

      expect(typeof image).toBe('string');
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
      expect(shape.attr('fill')).toContain(fill.attr('id'));
    });

    test('the image source is contributor avatar', async () => {
      const contributor = createMockContributor();
      const svg = await createContributorView(contributor, { size: 64 });

      expect(svg.findOne('pattern > image').attr('href')).toBe(contributor.avatar_url);
    });
  });
});
