import { Container, Element, Svg } from '@svgdotjs/svg.js';
import { createContributorAvatarImage, createSvgInstance } from './svg-builder';

describe('createSvgInstance()', () => {
  it('should return SVG instance', () => {
    const instance = createSvgInstance();
    expect(instance).toBeInstanceOf(Svg);
  });
});

describe('createRoundedImage', () => {
  let container: Container;

  beforeEach(() => {
    container = createSvgInstance();
  });

  it('should return svg element', () => {
    const el = createContributorAvatarImage(
      container,
      {
        login: 'foo',
        id: 1,
        html_url: 'a',
        avatar_url: 'b',
        contributions: 1,
      },
      64,
    );

    expect(el).toBeInstanceOf(Element);
  });

  it('should contain svg image', () => {
    const el = createContributorAvatarImage(
      container,
      {
        login: 'foo',
        id: 1,
        html_url: 'a',
        avatar_url: 'b',
        contributions: 1,
      },
      64,
    );

    expect(el.root().findOne('image').attr('href')).toBe('b');
  });

  it('should contain svg title', () => {
    const el = createContributorAvatarImage(
      container,
      {
        login: 'foo',
        id: 1,
        html_url: 'a',
        avatar_url: 'b',
        contributions: 1,
      },
      64,
    );

    expect(el.findOne('title').svg()).toContain('foo');
  });

  it('should contain svg link', () => {
    const el = createContributorAvatarImage(
      container,
      {
        login: 'foo',
        id: 1,
        html_url: 'a',
        avatar_url: 'b',
        contributions: 1,
      },
      64,
    );

    expect(el.root().findOne('a').attr('href')).toContain('a');
  });

  it('should be sized', () => {
    const el = createContributorAvatarImage(
      container,
      {
        login: 'foo',
        id: 1,
        html_url: 'a',
        avatar_url: 'b',
        contributions: 1,
      },
      64,
    );

    expect(el.width()).toBe(64);
    expect(el.height()).toBe(64);
  });
});
