export function validateRepoParam(repoParam: string | null) {
  if (!repoParam) {
    throw new Error(`'repo' parameter is required.`);
  }
  const [, repo] = repoParam.split('/');
  if (!repo) {
    throw new Error(`"${repoParam}" is not a repository string`);
  }
}
