import { Contributor, RendererOptions } from '@lib/core';
import * as SVG from '@svgdotjs/svg.js';
import { from, mergeMap, tap } from 'rxjs';
import { setupSvgRenderer } from './utils';
import * as rustRenderer from '@lib/renderer-rust';

setupSvgRenderer();

export type LayoutOptions = { itemSize: number; gap: number; maxColumns: number };
export interface ContributorsImageRenderer {
  readonly render: (contributors: Contributor[]) => Promise<string>;
  readonly layout: LayoutOptions;
}

export function createJsRenderer(options: RendererOptions): ContributorsImageRenderer {
  const itemSize = 64;
  const gap = 4;
  const maxColumns = options.maxColumns ?? 12;

  return {
    layout: { itemSize, gap, maxColumns },
    render: (contributors) => renderContributorsImage(contributors, { itemSize, gap, maxColumns }),
  };
}

export function createRustRenderer(options: RendererOptions): ContributorsImageRenderer {
  const itemSize = 64;
  const gap = 4;
  const maxColumns = options.maxColumns ?? 12;

  return {
    layout: { itemSize, gap, maxColumns },
    render: async (contributors) => rustRenderer.render(contributors, itemSize, maxColumns, gap),
  };
}

async function renderContributorsImage(contributors: Contributor[], layoutOptions: LayoutOptions): Promise<string> {
  const { itemSize, gap, maxColumns } = layoutOptions;
  const columns = Math.min(maxColumns, contributors.length);
  const rows = Math.ceil(contributors.length / columns);

  const container = await createContainer({ itemSize, gap, columns, rows });
  await renderContributors(container, contributors, { size: itemSize, columns, gap });

  return container.svg();
}

async function createContainer(params: {
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

async function renderContributors<T extends SVG.Element>(
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
