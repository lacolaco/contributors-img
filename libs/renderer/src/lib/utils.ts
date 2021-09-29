import * as SVG from '@svgdotjs/svg.js';
import AbortController from 'abort-controller';
import { JSDOM } from 'jsdom';
import { request } from 'undici';

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

    const { headers, body } = await request(imageUrl);
    const contentType = headers['content-type'];
    const buf = await body.arrayBuffer().then(Buffer.from);
    return `data:${contentType};base64,${buf.toString('base64')}`;
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
