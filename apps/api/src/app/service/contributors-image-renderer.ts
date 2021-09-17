import { Contributor } from '@lib/core';
import { renderContributorsImage } from '@lib/renderer';
import { injectable } from 'tsyringe';
import { runWithTracing } from '../utils/tracing';

@injectable()
export class ContributorsImageSvgRenderer {
  async renderSimpleAvatarTable(contributors: Contributor[]): Promise<string> {
    return runWithTracing('ContributorsImageSvgRenderer.renderSimpleAvatarTable', async () => {
      return renderContributorsImage(contributors);
    });
  }
}
