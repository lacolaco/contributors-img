import { Contributor } from '@lib/core';
import { renderContributorsImage } from '@lib/renderer';
import { injectable } from 'tsyringe';
import { ContentType } from '../utils/content-type';
import { runWithTracing } from '../utils/tracing';

@injectable()
export class ContributorsImageRenderer {
  async render(contributors: Contributor[]): Promise<{ data: string; contentType: ContentType }> {
    return runWithTracing('ContributorsImageRenderer.render', async () => {
      const svg = await renderContributorsImage(contributors);
      return {
        data: svg,
        contentType: ContentType.svg,
      };
    });
  }
}
