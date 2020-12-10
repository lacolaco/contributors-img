import { Repository } from './repository';

describe('Repository', () => {
  it('should create', () => {
    const repo = new Repository('foo', 'bar');
    expect(repo.owner).toBe('foo');
    expect(repo.repo).toBe('bar');
  });

  it('should create from repository string', () => {
    const repo = Repository.fromString('foo/bar');
    expect(repo.owner).toBe('foo');
    expect(repo.repo).toBe('bar');
  });

  it('should return a repository string', () => {
    const repo = new Repository('foo', 'bar');
    expect(repo.toString()).toBe('foo/bar');
  });
});
