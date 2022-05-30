import { Contributor, RendererOptions } from '@lib/core';
import * as rustRenderer from '@lib/renderer-rust';

export type LayoutOptions = { itemSize: number; gap: number; maxColumns: number };
export interface ContributorsImageRenderer {
  readonly render: (contributors: Contributor[]) => Promise<string>;
  readonly layout: LayoutOptions;
}

export function createRenderer(options: RendererOptions): ContributorsImageRenderer {
  const itemSize = 64;
  const gap = 4;
  const { maxColumns, maxCount } = options;

  return {
    layout: { itemSize, gap, maxColumns },
    render: async (contributors) => rustRenderer.render(contributors.slice(0, maxCount), itemSize, maxColumns, gap),
  };
}
