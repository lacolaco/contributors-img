import { TestBed } from '@angular/core/testing';

import { ContributorsService } from './contributors.service';

describe('ContributorsService', () => {
  beforeEach(() => TestBed.configureTestingModule({}));

  it('should be created', () => {
    const service: ContributorsService = TestBed.get(ContributorsService);
    expect(service).toBeTruthy();
  });
});
