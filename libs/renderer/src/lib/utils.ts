import * as SVG from '@svgdotjs/svg.js';
import AbortController from 'abort-controller';
import { JSDOM } from 'jsdom';
import fetch from 'node-fetch';

// Accept window from JSDOM
declare module '@svgdotjs/svg.js' {
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  export function registerWindow(windowImpl: any, documentImpl: any): void;
}

export async function createImageDataURI(imageUrl: string): Promise<string> {
  try {
    const controller = new AbortController();
    setTimeout(() => {
      controller.abort();
    }, 30 * 1000);
    const resp = await fetch(imageUrl, {
      signal: controller.signal,
    });
    const contentType = resp.headers.get('content-type');
    const data = (await resp.buffer()).toString('base64');
    return `data:${contentType};base64,${data}`;
  } catch (error) {
    console.error(error);
    return '';
  }
}

export function setupSvgRenderer(): void {
  const dom = new JSDOM();
  SVG.registerWindow(dom.window, dom.window.document);
}

export function createSvgInstance() {
  return SVG.SVG();
}
