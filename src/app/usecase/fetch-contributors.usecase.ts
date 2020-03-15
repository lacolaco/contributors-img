import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { GitHubContributor } from '../core/models';
import { ContributorsStore } from '../state/contributors';

@Injectable({
  providedIn: 'root',
})
export class FetchContributorsUsecase {
  constructor(private http: HttpClient, private store: ContributorsStore) {}

  async fetchContributors(repository: string) {
    this.store.update(state => ({
      ...state,
      contributors: {
        items: [],
        fetching: state.contributors.fetching + 1,
      },
      imageSnippet: '',
    }));

    try {
      const contributors = await this.http
        .get<GitHubContributor[]>(`/api/contributors`, { params: { repo: repository } })
        .toPromise();
      this.store.update(state => ({
        ...state,
        contributors: {
          items: contributors,
          fetching: state.contributors.fetching - 1,
        },
        imageSnippet: this.getImageSnippet(repository),
      }));
    } catch (error) {
      this.store.update(state => ({
        ...state,
        contributors: {
          items: [],
          fetching: state.contributors.fetching - 1,
        },
      }));
      throw error;
    }
  }

  private getImageSnippet(repository: string) {
    return `
<a href="https://github.com/${repository}/graphs/contributors">
  <img src="${location.origin}/image?repo=${repository}" />
</a>

Made with [contributors-img](${location.origin}).
`.trim();
  }
}
