import { Contributor } from '@lib/core';
import { injectable } from 'tsyringe';
import { createContributorAvatarImage, createSvgInstance } from '../utils/svg-builder';

@injectable()
export class ContributorsImageSvgRenderer {
  async renderSimpleAvatarTable(contributors: Contributor[]): Promise<string> {
    const avatarSize = 64;
    const gap = 4;
    const columnCount = Math.min(12, contributors.length);
    const rowCount = Math.ceil(contributors.length / columnCount);

    const container = createSvgInstance();
    container
      .size((avatarSize + gap) * (columnCount - 1) + avatarSize, (avatarSize + gap) * (rowCount - 1) + avatarSize)
      .css('padding', '4px');

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
}
