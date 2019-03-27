import { Injectable } from '@angular/core';
import { Store } from '@lacolaco/ngx-store';
import { GitHubContributor } from '../core/models';

export interface State {
  repository: string;
  items: GitHubContributor[];
  itemsLoading: number;
  imageSnippet: string;
}

export const initialState: State = {
  repository: 'angular/angular-ja',
  items: [],
  itemsLoading: 0,
  imageSnippet: '',
};

@Injectable({ providedIn: 'root' })
export class ContributorsStore extends Store<State> {
  constructor() {
    super({ initialState });
    // TODO: use DI
    const repoFromUrl = new URLSearchParams(window.location.search).get('repo');
    if (repoFromUrl && repoFromUrl.trim().length > 0) {
      this.updateValue(state => ({ ...state, repository: repoFromUrl }));
    }
  }
}
