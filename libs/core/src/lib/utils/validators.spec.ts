import { assertRepositoryName } from './validators';

describe('validators', () => {
  describe('validateRepoParam', () => {
    test('should accept owner/repo string', () => {
      expect(assertRepositoryName('foo/bar')).toBe(true);
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
