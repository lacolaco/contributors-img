import { registerWindow, SVG } from '@svgdotjs/svg.js';
import { createWindow } from 'domino';

declare module '@svgdotjs/svg.js' {
  export function registerWindow(windowImpl: Window, documentImpl: Document): void;
}

const windowImpl = createWindow();
registerWindow(windowImpl, windowImpl.document);

export function createSvgInstance() {
  return SVG();
}
