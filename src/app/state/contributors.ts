import { Injectable } from '@angular/core';
import { Store } from '@lacolaco/reactive-store';
import { GitHubContributor } from '../core/models';

export interface State {
  repository: string;
  contributors: {
    items: GitHubContributor[];
    fetching: number;
  };
  imageSnippet: string;
}

export const initialValue: State = {
  repository: 'angular/angular-ja',
  contributors: {
    items: [],
    fetching: 0,
  },
  imageSnippet: '',
};

@Injectable({ providedIn: 'root' })
export class ContributorsStore extends Store<State> {
  constructor() {
    super({ initialValue });
    // TODO: use DI
    const repoFromUrl = new URLSearchParams(window.location.search).get('repo');
    if (repoFromUrl && repoFromUrl.trim().length > 0) {
      this.update(state => ({
        ...state,
        repository: repoFromUrl,
      }));
    }
  }
}
