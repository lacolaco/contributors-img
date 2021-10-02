import { Contributor } from '@lib/core';
import { createRenderer } from '@lib/renderer';
import { from } from 'rxjs';
import { concatMap, toArray } from 'rxjs/operators';
import { injectable } from 'tsyringe';
import { request } from 'undici';
import { ContentType } from '../utils/content-type';
import { runWithTracing } from '../utils/tracing';

const convertImageToDataURL = async (imageURL: string, imageSize: number): Promise<string> => {
  const url = new URL(imageURL);
  url.searchParams.set('size', imageSize.toString());
  try {
    const { headers, body } = await request(url, { method: 'GET' });
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
      const renderer = createRenderer();

      const svg = await from(contributors)
        .pipe(
          concatMap(async (contributor) => ({
            ...contributor,
            avatar_url: await convertImageToDataURL(contributor.avatar_url, renderer.layout.itemSize),
          })),
          toArray(),
          concatMap((contributors) => renderer.render(contributors)),
        )
        .toPromise();

      return {
        data: svg,
        contentType: ContentType.svg,
      };
    });
  }
}
