import { Octokit } from '@octokit/rest';
import { environment } from '../../../environments/environment';

export const createOctokit = (): Octokit => {
  return new Octokit({
    auth: environment.githubAuthToken,
  });
};
