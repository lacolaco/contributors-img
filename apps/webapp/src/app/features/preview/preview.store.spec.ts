import { Repository } from '@lib/core';
import { createServiceFactory, SpectatorService } from '@ngneat/spectator/jest';
import { PreviewStore } from './preview.store';

describe('AppStore', () => {
  let spectator: SpectatorService<PreviewStore>;
  const createService = createServiceFactory(PreviewStore);

  beforeEach(() => (spectator = createService()));

  it('should be created', () => {
    expect(spectator.service).toBeTruthy();
  });

  it('should has initial value', () => {
    expect(spectator.service.value).toBeDefined();
    expect(spectator.service.value.repository).toBeNull();
    expect(spectator.service.value.image.data).toEqual(null);
    expect(spectator.service.value.image.fetching).toEqual(0);
  });

  describe('startFetchingContributors()', () => {
    it('should update value', () => {
      const repo = new Repository('foo', 'bar');
      spectator.service.startFetchingImage(repo);

      expect(spectator.service.value.repository).toBe(repo);
      expect(spectator.service.value.image.data).toEqual(null);
      expect(spectator.service.value.image.fetching).toEqual(1);
    });
  });

  describe('finishFetchingContributors()', () => {
    it('should update value', () => {
      spectator.service.update((state) => ({
        ...state,
        image: {
          data: null,
          fetching: 1,
        },
      }));

      spectator.service.finishFetchingImage();

      expect(spectator.service.value.image.fetching).toEqual(0);
    });
  });
});
