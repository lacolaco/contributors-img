import { Contributor } from '../../index';
import faker from '@faker-js/faker';

export function createMockContributor(params: Partial<Contributor> = {}): Contributor {
  const username = faker.internet.userName();
  return {
    id: faker.datatype.number({ min: 0 }),
    login: username,
    avatar_url: faker.internet.avatar(),
    html_url: `https://github.com/${username}`,
    contributions: faker.datatype.number({ min: 0 }),
    ...params,
  };
}
