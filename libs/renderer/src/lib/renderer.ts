import { Contributor } from '@lib/core';
import * as SVG from '@svgdotjs/svg.js';
import { from } from 'rxjs';
import { mergeMap, tap } from 'rxjs/operators';
import { createImageDataURI, setupSvgRenderer } from './utils';

setupSvgRenderer();

export async function renderContributorsImage(contributors: Contributor[]): Promise<string> {
  const itemSize = 64;
  const gap = 4;
  const columns = Math.min(12, contributors.length);
  const rows = Math.ceil(contributors.length / columns);

  const container = await createContainer({ itemSize, gap, columns, rows });
  await renderContributors(container, contributors, { size: itemSize, columns, gap });
  await inlineForeignImages(container);

  return container.svg();
}

export async function createContainer(params: {
  itemSize: number;
  columns: number;
  rows: number;
  gap: number;
}): Promise<SVG.Container> {
  const { itemSize, columns, rows, gap } = params;
  const width = itemSize * columns + gap * (columns - 1);
  const height = itemSize * rows + gap * (rows - 1);

  const container = SVG.SVG();
  container.size(width, height);
  return container;
}

export async function renderContributors<T extends SVG.Element>(
  container: T,
  contributors: Contributor[],
  params: { size: number; columns: number; gap: number },
): Promise<void> {
  const { size, columns, gap } = params;
  await from(contributors.entries())
    .pipe(
      mergeMap(async ([i, contributor]) => {
        return [i, await createContributorView(contributor, { size })] as const;
      }),
      tap(([i, image]) => {
        const x = (i % columns) * (size + gap);
        const y = Math.floor(i / columns) * (size + gap);
        image.addTo(container).move(x, y);
      }),
    )
    .toPromise();
}

export async function inlineForeignImages(element: SVG.Element): Promise<void> {
  return await from(element.find('image'))
    .pipe(
      mergeMap(async (image) => {
        const href = image.attr('href') as string;
        if (href.startsWith('data:')) {
          return;
        }
        const imageURL = new URL(href);
        imageURL.searchParams.set('size', image.width().toString());
        const uri = await createImageDataURI(imageURL.toString());
        image.attr('href', uri);
      }),
    )
    .toPromise();
}

export async function createContributorView(contributor: Contributor, options: { size: number }) {
  const { size } = options;
  const borderColor = '#c0c0c0';

  const view = SVG.SVG();
  view.size(size, size);

  // add title
  const title = SVG.SVG(`<title>${contributor.login}</title>`);
  view.add(title);

  // add a shape
  const borderWidth = 1;
  const shape = new SVG.Circle();
  shape
    .size(size - borderWidth)
    .center(size / 2, size / 2)
    .stroke({ width: borderWidth, color: borderColor });
  view.add(shape);

  // fill shape with image
  const imageFill = view.pattern(size, size);
  imageFill.image(contributor.avatar_url).size(size, size);
  shape.fill(imageFill);

  return view;
}
