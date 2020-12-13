import { Octokit } from '@octokit/rest';

export const createOctokit = (): Octokit => {
  return new Octokit({
    auth: 'token 393ad1f410e7f6e6d78a19466812b6cea4d1ed52',
  });
};
