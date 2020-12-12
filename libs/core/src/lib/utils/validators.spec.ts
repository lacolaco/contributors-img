import { assertRepositoryName } from './validators';

describe('validators', () => {
  describe('validateRepoParam', () => {
    it('should accept owner/repo string', () => {
      expect(() => assertRepositoryName('foo/bar')).not.toThrow();
    });

    it('should reject empty string', () => {
      expect(() => assertRepositoryName('')).toThrow();
    });

    it('should reject null', () => {
      expect(() => assertRepositoryName(null)).toThrow();
    });

    it('should reject non repository string', () => {
      expect(() => assertRepositoryName('foo')).toThrow();
    });

    it('should reject non repository string', () => {
      expect(() => assertRepositoryName('foo/')).toThrow();
    });
  });
});
