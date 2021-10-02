import { Contributor } from '@lib/core';
import { renderContributorsImage } from '@lib/renderer';
import AbortController from 'abort-controller';
import { injectable } from 'tsyringe';
import { request } from 'undici';
import { ContentType } from '../utils/content-type';
import { runWithTracing } from '../utils/tracing';

const imageUriTransformer = async (imageUrl: string): Promise<string> => {
  try {
    const controller = new AbortController();
    setTimeout(() => {
      controller.abort();
    }, 30 * 1000);

    const { headers, body } = await request(imageUrl, {
      method: 'GET',
      signal: controller.signal,
    });
    const contentType = headers['content-type'];
    const buf = await body.arrayBuffer().then(Buffer.from);
    return `data:${contentType};base64,${buf.toString('base64')}`;
  } catch (error) {
    console.error(error);
    return '';
  }
};

@injectable()
export class ContributorsImageRenderer {
  async render(contributors: Contributor[]): Promise<{ data: string; contentType: ContentType }> {
    return runWithTracing('ContributorsImageRenderer.render', async () => {
      const svg = await renderContributorsImage(contributors, imageUriTransformer);
      return {
        data: svg,
        contentType: ContentType.svg,
      };
    });
  }
}
