import { ContributorsImageCache } from '../service/image-cache';
import { renderContributorsImage } from '../service/render-image';
import { Contributor, Repository } from '../shared/model';

export function createContributorsImageQuery(
  imageCache: ContributorsImageCache,
  config: { webappUrl: string; useHeadless: boolean },
): (repository: Repository, contributors: Contributor[]) => Promise<Buffer> {
  return async (repository, contributors) => {
    console.debug('restore cache');
    const cached = await imageCache.restore(repository);
    if (cached) {
      console.debug('cache hit');
      return cached;
    }
    console.debug(`render image`);
    const rendered = await renderContributorsImage({
      contributors,
      config: {
        webappUrl: config.webappUrl,
        useHeadless: config.useHeadless,
      },
    });
    console.debug('save cache');
    await imageCache.save(repository, rendered);
    return rendered;
  };
}
