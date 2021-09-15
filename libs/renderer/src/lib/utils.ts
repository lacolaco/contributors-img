import * as SVG from '@svgdotjs/svg.js';
import AbortController from 'abort-controller';
import { JSDOM } from 'jsdom';
import fetch from 'node-fetch';

// Accept window from JSDOM
declare module '@svgdotjs/svg.js' {
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  export function registerWindow(windowImpl: any, documentImpl: any): void;
}

export async function createDataURIFromURL(imageUrl: string): Promise<string> {
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

export function createSvgInstance() {
  if (SVG.getWindow() == null) {
    const DOM = new JSDOM();
    SVG.registerWindow(DOM.window, DOM.window.document);
  }
  return SVG.SVG();
}
