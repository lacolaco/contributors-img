import { Repository } from '@lib/core';
import { PreviewStore } from './preview.store';

describe('PreviewStore', () => {
  let store: PreviewStore;

  beforeEach(() => (store = new PreviewStore()));

  it('should be created', () => {
    expect(store).toBeTruthy();
  });

  it('should has initial value', () => {
    expect(store.value).toBeDefined();
    expect(store.value.repository).toBeNull();
    expect(store.value.image.data).toEqual(null);
    expect(store.value.image.fetching).toEqual(0);
  });

  describe('startFetchingContributors()', () => {
    it('should update value', () => {
      const repo = new Repository('foo', 'bar');
      store.startFetchingImage(repo);

      expect(store.value.repository).toBe(repo);
      expect(store.value.image.data).toEqual(null);
      expect(store.value.image.fetching).toEqual(1);
    });
  });

  describe('finishFetchingContributors()', () => {
    it('should update value', () => {
      store.update((state) => ({
        ...state,
        image: {
          data: null,
          fetching: 1,
        },
      }));

      store.finishFetchingImage();

      expect(store.value.image.fetching).toEqual(0);
    });
  });
});
