import * as SVG from '@svgdotjs/svg.js';
import { JSDOM } from 'jsdom';

declare module '@svgdotjs/svg.js' {
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  export function registerWindow(windowImpl: any, documentImpl: any): void;
}

const DOM = new JSDOM();
SVG.registerWindow(DOM.window, DOM.window.document);

export function createSvgInstance() {
  return SVG.SVG();
}

