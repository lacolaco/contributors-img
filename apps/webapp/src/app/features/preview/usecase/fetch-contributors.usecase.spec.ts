import { createServiceFactory, SpectatorService } from '@ngneat/spectator/jest';
import { FetchContributorsUsecase } from './fetch-contributors.usecase';
import { HttpClientTestingModule } from '@angular/common/http/testing';

describe('FetchContributorsUsecase', () => {
  let spectator: SpectatorService<FetchContributorsUsecase>;
  const createService = createServiceFactory({
    service: FetchContributorsUsecase,
    imports: [HttpClientTestingModule],
  });

  beforeEach(() => (spectator = createService()));

  test('should be created', () => {
    expect(spectator.service).toBeTruthy();
  });
});
