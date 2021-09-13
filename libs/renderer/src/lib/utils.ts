import { Contributor } from '@lib/core';
import * as SVG from '@svgdotjs/svg.js';
import { JSDOM } from 'jsdom';

import fetch from 'node-fetch';
import AbortController from 'abort-controller';

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

declare module '@svgdotjs/svg.js' {
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  export function registerWindow(windowImpl: any, documentImpl: any): void;
}

const DOM = new JSDOM();
SVG.registerWindow(DOM.window, DOM.window.document);

export function createSvgInstance() {
  return SVG.SVG();
}

export async function createContributorAvatarImage(
  container: SVG.Container,
  { login, avatar_url, html_url }: Contributor,
  size: number,
  borderColor = '#c0c0c0',
) {
  const imageURL = new URL(avatar_url);
  imageURL.searchParams.set('size', size.toString());
  const imageURI = await createDataURIFromURL(imageURL.toString());

  const image = container.image(imageURI).size(size, size);
  const bg = container.pattern(size, size).add(image);
  const link = container.link(html_url).target('_blank');
  return container
    .circle(size, size)
    .stroke({ color: borderColor, width: 1 })
    .fill(bg)
    .add(container.element('title').words(login))
    .addTo(link);
}
