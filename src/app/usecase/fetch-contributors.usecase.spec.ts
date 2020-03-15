import { createServiceFactory, SpectatorService } from '@ngneat/spectator';
import { FetchContributorsUsecase } from './fetch-contributors.usecase';

describe('FetchContributorsUsecase', () => {
  let spectator: SpectatorService<FetchContributorsUsecase>;
  const createService = createServiceFactory(FetchContributorsUsecase);

  beforeEach(() => (spectator = createService()));

  it('should be created', () => {
    expect(spectator.service).toBeTruthy();
  });
});
