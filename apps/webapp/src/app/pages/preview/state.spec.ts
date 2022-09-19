import { defaultImageParams, PreviewState } from './state';

describe('PreviewState', () => {
  let state: PreviewState;

  beforeEach(() => (state = new PreviewState()));

  it('should be created', () => {
    expect(state).toBeTruthy();
  });

  it('should has initial value', () => {
    expect(state.get()).toBeDefined();
    expect(state.get().imageParams).toBe(defaultImageParams);
    expect(state.get().result).toEqual(null);
    expect(state.get().fetchingCount).toEqual(0);
  });

  describe('startFetchingContributors()', () => {
    it('should update value', () => {
      state.startFetchingImage();

      expect(state.get().result).toEqual(null);
      expect(state.get().fetchingCount).toEqual(1);
    });
  });

  describe('finishFetchingContributors()', () => {
    it('should update value', () => {
      state.set((state) => ({
        ...state,
        fetchingCount: 1,
        result: null,
      }));

      state.finishFetchingImage(null);

      expect(state.get().fetchingCount).toEqual(0);
    });
  });
});
