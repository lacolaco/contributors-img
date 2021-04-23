import { Contributor } from '@lib/core';
import * as SVG from '@svgdotjs/svg.js';
import { JSDOM } from 'jsdom';
import { createDataURIFromURL } from './image';

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
