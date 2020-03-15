import { ContributorsStore } from './contributors';
import { SpectatorService, createServiceFactory } from '@ngneat/spectator';

describe('ContributorsStore', () => {
  let spectator: SpectatorService<ContributorsStore>;
  const createService = createServiceFactory(ContributorsStore);

  beforeEach(() => (spectator = createService()));

  it('should be created', () => {
    expect(spectator.service).toBeTruthy();
  });

  it('should has initial value', () => {
    expect(spectator.service.value).toBeDefined();
    expect(spectator.service.value.repository).toBe('angular/angular-ja');
  });
});
