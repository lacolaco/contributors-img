import { Contributor } from '../../index';
import { datatype, internet } from 'faker';

export function createMockContributor(params: Partial<Contributor> = {}): Contributor {
  const username = internet.userName();
  return {
    id: datatype.number({ min: 0 }),
    login: username,
    avatar_url: internet.avatar(),
    html_url: `https://github.com/${username}`,
    contributions: datatype.number({ min: 0 }),
    ...params,
  };
}
