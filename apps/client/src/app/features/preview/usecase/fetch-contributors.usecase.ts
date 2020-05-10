import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Contributor, Repository } from '@contributors-img/api-interfaces';
import { PreviewStore } from '../store';

@Injectable({
  providedIn: 'root',
})
export class FetchContributorsUsecase {
  constructor(private http: HttpClient, private store: PreviewStore) {}

  async execute(repoName: string) {
    this.store.startFetchingContributors(Repository.fromString(repoName));

    try {
      const contributors = await this.http
        .get<Contributor[]>(`/api/contributors`, { params: { repo: repoName } })
        .toPromise();
      this.store.finishFetchingContributors(contributors);
    } catch (error) {
      console.error(error);
      this.store.finishFetchingContributors([]);
    }
  }
}
