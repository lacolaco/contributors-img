import { Contributor } from '@lib/core';
import * as SVG from '@svgdotjs/svg.js';
import { createDataURIFromURL, createSvgInstance } from './utils';

export async function renderContributorsImage(contributors: Contributor[]): Promise<string> {
  const avatarSize = 64;
  const gap = 4;
  const columnCount = Math.min(12, contributors.length);
  const rowCount = Math.ceil(contributors.length / columnCount);

  const container = createSvgInstance();
  container.size((avatarSize + gap) * (columnCount - 1) + avatarSize, (avatarSize + gap) * (rowCount - 1) + avatarSize);

  await Promise.all(
    Array.from(contributors.entries()).map(([i, contributor]) => {
      const x = (i % columnCount) * (avatarSize + gap);
      const y = Math.floor(i / columnCount) * (avatarSize + gap);

      const inner = container.nested().move(x, y);
      return createContributorAvatarImage(inner, contributor, avatarSize);
    }),
  );

  return container.svg();
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
