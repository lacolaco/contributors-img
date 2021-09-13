import { Container, Element, Svg } from '@svgdotjs/svg.js';
import { createContributorAvatarImage, createSvgInstance } from './utils';

describe('createSvgInstance()', () => {
  it('should return SVG instance', () => {
    const instance = createSvgInstance();
    expect(instance).toBeInstanceOf(Svg);
  });
});

describe('createRoundedImage', () => {
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
