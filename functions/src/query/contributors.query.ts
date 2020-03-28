import { Repository, Contributor } from '../shared/model';
import { ContributorsJsonCache } from '../service/json-cache';
import { fetchContributors } from '../service/fetch-contributors';

export function createContributorsQuery(
  cache: ContributorsJsonCache,
): (repository: Repository) => Promise<Contributor[]> {
  return async (repository) => {
    console.debug('restore cache');
    const cached = await cache.restore(repository);
    if (cached) {
      console.debug('cache hit');
      return cached;
    }
    console.debug(`fetch data`);
    const contributors = await fetchContributors(repository);
    console.debug('save cache');
    await cache.save(repository, contributors);
    return contributors;
  };
}
