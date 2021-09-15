import { Svg } from '@svgdotjs/svg.js';
import { createSvgInstance, setupSvgRenderer } from './utils';

describe('createSvgInstance()', () => {
  it('should return SVG instance', () => {
    setupSvgRenderer();
    const instance = createSvgInstance();
    expect(instance).toBeInstanceOf(Svg);
  });
});
