import { Injectable } from '@angular/core';
import { Store } from '@lacolaco/reactive-store';
import { GitHubContributor } from '../core/models';
import { Repository } from '@api/shared/model/repository';

export interface State {
  repository: Repository | null;
  contributors: {
    items: GitHubContributor[];
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
export class AppStore extends Store<State> {
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

  finishFetchingContributors(items: GitHubContributor[]) {
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
