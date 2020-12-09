import { Octokit } from '@octokit/rest';
import { Contributor } from '../shared/model/contributor';
import { Repository } from '../shared/model/repository';

const octokit = new Octokit({
  auth: 'token 393ad1f410e7f6e6d78a19466812b6cea4d1ed52',
});

export async function fetchContributors(repo: Repository): Promise<Contributor[]> {
  // Fetch all contributors with auto-pagination
  return await octokit.paginate(octokit.repos.listContributors, {
    owner: repo.owner,
    repo: repo.repo,
    per_page: 100,
  });
}
