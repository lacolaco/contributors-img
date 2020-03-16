import { createServiceFactory, SpectatorService } from '@ngneat/spectator';
import { Repository } from 'shared/model/repository';
import { AppStore } from './store';

describe('AppStore', () => {
  let spectator: SpectatorService<AppStore>;
  const createService = createServiceFactory(AppStore);

  beforeEach(() => (spectator = createService()));

  it('should be created', () => {
    expect(spectator.service).toBeTruthy();
  });

  it('should has initial value', () => {
    expect(spectator.service.value).toBeDefined();
    expect(spectator.service.value.repository).toBeNull();
    expect(spectator.service.value.contributors.items).toEqual([]);
    expect(spectator.service.value.contributors.fetching).toEqual(0);
  });

  describe('startFetchingContributors()', () => {
    it('should update value', () => {
      const repo = new Repository('foo', 'bar');
      spectator.service.startFetchingContributors(repo);

      expect(spectator.service.value.repository).toBe(repo);
      expect(spectator.service.value.contributors.items).toEqual([]);
      expect(spectator.service.value.contributors.fetching).toEqual(1);
    });
  });

  describe('finishFetchingContributors()', () => {
    it('should update value', () => {
      spectator.service.update(state => ({
        ...state,
        contributors: {
          items: [],
          fetching: 1,
        },
      }));

      spectator.service.finishFetchingContributors([
        { id: 1, login: 'foo', avatar_url: '', html_url: '', contributions: 1 },
      ]);

      expect(spectator.service.value.contributors.items.length).toEqual(1);
      expect(spectator.service.value.contributors.fetching).toEqual(0);
    });
  });
});
