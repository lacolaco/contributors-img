import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { GitHubContributor } from '../core/models';
import { AppStore } from '../state/store';
import { Repository } from '@api/shared/model/repository';

@Injectable({
  providedIn: 'root',
})
export class FetchContributorsUsecase {
  constructor(private http: HttpClient, private store: AppStore) {}

  async fetchContributors(repoName: string) {
    this.store.startFetchingContributors(Repository.fromString(repoName));

    try {
      const contributors = await this.http
        .get<GitHubContributor[]>(`/api/contributors`, { params: { repo: repoName } })
        .toPromise();
      this.store.finishFetchingContributors(contributors);
    } catch (error) {
      console.error(error);
      this.store.finishFetchingContributors([]);
    }
  }
}
