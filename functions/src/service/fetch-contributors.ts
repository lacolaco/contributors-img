import { Octokit } from '@octokit/rest';
import { Repository } from '../model/repository';

const octokit = new Octokit({
  auth: 'token 393ad1f410e7f6e6d78a19466812b6cea4d1ed52',
});

export async function fetchContributors(repo: Repository): Promise<Octokit.ReposListContributorsResponse> {
  // Fetch all contributors with auto-pagination
  const options = octokit.repos.listContributors.endpoint.merge({
    owner: repo.owner,
    repo: repo.repo,
    per_page: 100,
  });
  return await octokit.paginate(options);
}
