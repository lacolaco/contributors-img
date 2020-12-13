export function assertRepositoryName(repoName: string | unknown): repoName is string {
  if (typeof repoName !== 'string') {
    return false;
  }
  if (!/^[\w\-._]+\/[\w\-._]+$/.test(repoName)) {
    return false;
  }
  return true;
}
