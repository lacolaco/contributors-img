export function validateRepoParam(repoParam: unknown | null): asserts repoParam is string {
  if (!repoParam || typeof repoParam !== 'string') {
    throw new Error(`'repo' parameter is required.`);
  }
  const [, repo] = repoParam.split('/');
  if (!repo) {
    throw new Error(`"${repoParam}" is not a repository string`);
  }
}
