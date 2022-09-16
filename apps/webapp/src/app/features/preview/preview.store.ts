import { Store } from '@lacolaco/reactive-store';
import { Repository } from '../../shared/model/repository';

export interface State {
  repository: Repository | null;
  image: {
    data: string | null;
    fetching: number;
  };
  showImageSnippet: boolean;
}

export const initialValue: State = {
  repository: null,
  image: {
    data: null,
    fetching: 0,
  },
  showImageSnippet: false,
};

export class PreviewStore extends Store<State> {
  constructor() {
    super({ initialValue });
  }

  startFetchingImage(repository: Repository) {
    this.update((state) => ({
      ...state,
      repository,
      image: {
        ...state.image,
        data: null,
        fetching: state.image.fetching + 1,
      },
    }));
  }

  finishFetchingImage() {
    this.update((state) => ({
      ...state,
      image: {
        ...state.image,
        fetching: state.image.fetching - 1,
      },
    }));
  }

  setImageData(data: string) {
    this.update((state) => ({
      ...state,
      image: {
        ...state.image,
        data,
      },
    }));
  }

  showImageSnippet() {
    this.update((state) => ({
      ...state,
      showImageSnippet: true,
    }));
  }

  closeImageSnippet() {
    this.update((state) => ({
      ...state,
      showImageSnippet: false,
    }));
  }
}
