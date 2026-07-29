import { TestBed } from '@angular/core/testing';
import { defaultImageParams, PreviewState } from './state';

describe('PreviewState', () => {
  let state: PreviewState;

  beforeEach(() => {
    TestBed.runInInjectionContext(() => {
      state = new PreviewState();
    });
  });

  it('should be created', () => {
    expect(state).toBeTruthy();
  });

  it('should has initial value', () => {
    expect(state.state()).toBeDefined();
    expect(state.state().imageParams).toBe(defaultImageParams);
    expect(state.state().result).toEqual(null);
    expect(state.state().fetchingCount).toEqual(0);
  });

  describe('startFetchingContributors()', () => {
    it('should update value', () => {
      state.startFetchingImage();

      expect(state.state().result).toEqual(null);
      expect(state.state().fetchingCount).toEqual(1);
    });
  });

  describe('finishFetchingContributors()', () => {
    it('should update value', () => {
      state.startFetchingImage();

      state.finishFetchingImage(null);

      expect(state.state().fetchingCount).toEqual(0);
    });
  });
});
