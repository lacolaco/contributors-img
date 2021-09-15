import { renderContributorsImage, createContributorAvatarImage } from './renderer';
import { Contributor } from '@lib/core';
import { Container, Element } from '@svgdotjs/svg.js';
import { createSvgInstance } from './utils';

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

describe('createContributorAvatarImage', () => {
  let container: Container;

  const createContributor = () => ({
    login: 'lacolaco',
    id: 1,
    html_url: 'https://github.com/lacolaco',
    avatar_url: 'https://github.com/lacolaco.png',
    contributions: 1,
  });

  beforeEach(() => {
    container = createSvgInstance();
  });

  it('should return svg element', async () => {
    const el = await createContributorAvatarImage(container, createContributor(), 64);

    expect(el).toBeInstanceOf(Element);
  });

  it('should contain svg image', async () => {
    const el = await createContributorAvatarImage(container, createContributor(), 64);

    expect(el.root().findOne('image')).toBeTruthy();
  });

  it('should contain svg title', async () => {
    const contributor = createContributor();
    const el = await createContributorAvatarImage(container, contributor, 64);

    expect(el.findOne('title').svg()).toContain(contributor.login);
  });

  it('should contain svg link', async () => {
    const contributor = createContributor();
    const el = await createContributorAvatarImage(container, contributor, 64);

    expect(el.root().findOne('a').attr('href')).toBe(contributor.html_url);
  });

  it('should be sized', async () => {
    const el = await createContributorAvatarImage(container, createContributor(), 64);

    expect(el.width()).toBe(64);
    expect(el.height()).toBe(64);
  });
});
