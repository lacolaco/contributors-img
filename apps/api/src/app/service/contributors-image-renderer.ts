import { Contributor } from '@lib/core';
import { renderContributorsImage } from '@lib/renderer';
import { injectable } from 'tsyringe';

@injectable()
export class ContributorsImageSvgRenderer {
  async renderSimpleAvatarTable(contributors: Contributor[]): Promise<string> {
    return renderContributorsImage(contributors);
  }
}
