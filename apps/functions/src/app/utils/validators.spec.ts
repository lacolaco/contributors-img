import { validateRepoParam } from './validators';

describe('validators', () => {
  describe('validateRepoParam', () => {
    it('should accept owner/repo string', () => {
      expect(() => validateRepoParam('foo/bar')).not.toThrow();
    });

    it('should reject empty string', () => {
      expect(() => validateRepoParam('')).toThrow();
    });

    it('should reject null', () => {
      expect(() => validateRepoParam(null)).toThrow();
    });

    it('should reject non repository string', () => {
      expect(() => validateRepoParam('foo')).toThrow();
    });

    it('should reject non repository string', () => {
      expect(() => validateRepoParam('foo/')).toThrow();
    });
  });
});
