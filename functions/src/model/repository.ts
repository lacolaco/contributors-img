export class Repository {
  constructor(public owner: string, public repo: string) {}

  static fromString(repoStr: string): Repository {
    const [owner, repo] = repoStr.split('/');
    return new Repository(owner, repo);
  }

  toString(): string {
    return `${this.owner}/${this.repo}`;
  }
}
