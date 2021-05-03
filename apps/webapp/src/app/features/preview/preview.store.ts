import { Injectable } from '@angular/core';
import { Store } from '@lacolaco/reactive-store';
import { Repository } from '@lib/core';

export interface State {
  repository: Repository | null;
  image: {
    data: string | null;
    fetching: number;
  };
}

export const initialValue: State = {
  repository: null,
  image: {
    data: null,
    fetching: 0,
  },
};

@Injectable({ providedIn: 'root' })
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

  finishFetchingImage(data: string) {
    this.update((state) => ({
      ...state,
      image: {
        ...state.image,
        data,
        fetching: state.image.fetching - 1,
      },
    }));
  }
}
