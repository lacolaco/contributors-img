import { Contributor } from './contributor';

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

export interface RepositoryContributors {
  owner: string;
  repo: string;
  stargazersCount: number;
  data: Contributor[];
}

export interface FeaturedRepository {
  repository: string;
  stargazers: number;
  contributors: number;
}
