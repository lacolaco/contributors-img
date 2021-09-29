import * as SVG from '@svgdotjs/svg.js';
import { JSDOM } from 'jsdom';

// Accept window from JSDOM
declare module '@svgdotjs/svg.js' {
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  export function registerWindow(windowImpl: any, documentImpl: any): void;
}

export function setupSvgRenderer(): void {
  const dom = new JSDOM();
  SVG.registerWindow(dom.window, dom.window.document);
}

export function createSvgInstance() {
  return SVG.SVG();
}
