import { Svg } from '@svgdotjs/svg.js';
import { createSvgInstance } from './utils';

describe('createSvgInstance()', () => {
  it('should return SVG instance', () => {
    const instance = createSvgInstance();
    expect(instance).toBeInstanceOf(Svg);
  });
});
