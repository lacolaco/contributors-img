import { assertRepositoryName } from './validators';

describe('validators', () => {
  describe('validateRepoParam', () => {
    test('should accept owner/repo string', () => {
      expect(assertRepositoryName('foo/bar')).toBe(true);
    });

    test('should accept hyphens in repo name', () => {
      expect(assertRepositoryName('foo-org/bar-project')).toBe(true);
    });

    test('should accept underscores in repo name', () => {
      expect(assertRepositoryName('foo_org/bar_project')).toBe(true);
    });

    test('should accept periods in repo name', () => {
      expect(assertRepositoryName('foo.org/bar.js')).toBe(true);
    });

    test('should reject empty string', () => {
      expect(assertRepositoryName('')).toBe(false);
    });

    test('should reject null', () => {
      expect(assertRepositoryName(null)).toBe(false);
    });

    test('should reject non repository string', () => {
      expect(assertRepositoryName('foo')).toBe(false);
    });

    test('should reject non repository string', () => {
      expect(assertRepositoryName('foo/')).toBe(false);
    });
  });
});
