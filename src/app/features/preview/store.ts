import { Injectable } from '@angular/core';
import { Contributor, Repository } from '@api/shared/model';
import { Store } from '@lacolaco/reactive-store';

export interface State {
  repository: Repository | null;
  contributors: {
    items: Contributor[];
    fetching: number;
  };
}

export const initialValue: State = {
  repository: null,
  contributors: {
    items: [],
    fetching: 0,
  },
};

@Injectable({ providedIn: 'root' })
export class PreviewStore extends Store<State> {
  constructor() {
    super({ initialValue });
  }

  startFetchingContributors(repository: Repository) {
    this.update(state => ({
      ...state,
      repository,
      contributors: {
        ...state.contributors,
        items: [],
        fetching: state.contributors.fetching + 1,
      },
    }));
  }

  finishFetchingContributors(items: Contributor[]) {
    this.update(state => ({
      ...state,
      contributors: {
        ...state.contributors,
        items,
        fetching: state.contributors.fetching - 1,
      },
    }));
  }
}
