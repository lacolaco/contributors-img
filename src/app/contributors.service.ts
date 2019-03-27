import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { environment } from 'src/environments/environment';
import { GitHubContributor } from './core/models';
import { ContributorsStore } from './state/contributors';

@Injectable({
  providedIn: 'root',
})
export class ContributorsService {
  constructor(private http: HttpClient, private store: ContributorsStore) {}

  async fetchContributors(repository: string) {
    this.store.updateValue(state => ({
      ...state,
      items: [],
      itemsLoading: state.itemsLoading + 1,
      imageSnippet: '',
    }));

    try {
      const contributors = await this.http
        .get<GitHubContributor[]>(`https://us-central1-contributors-img.cloudfunctions.net/getContributors`, {
          params: { repo: repository },
        })
        .toPromise();
      this.store.updateValue(state => ({
        ...state,
        items: contributors,
        imageSnippet: this.getImageSnippet(repository),
        itemsLoading: state.itemsLoading - 1,
      }));
    } catch (error) {
      this.store.updateValue(state => ({
        ...state,
        itemsLoading: state.itemsLoading - 1,
      }));
      throw error;
    }
  }

  private getImageSnippet(repository: string) {
    return `
<a href="https://github.com/${repository}/graphs/contributors">
  <img src="https://contributors-img.firebaseapp.com/image?repo=${repository}" />
</a>

Made with [contributors-img](https://contributors-img.firebaseapp.com).
`.trim();
  }
}
