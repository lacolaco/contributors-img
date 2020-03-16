import { Component, OnInit } from '@angular/core';
import { AppStore } from './state/store';
import { FetchContributorsUsecase } from './usecase/fetch-contributors.usecase';

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.scss'],
})
export class AppComponent implements OnInit {
  constructor(private usecase: FetchContributorsUsecase, private store: AppStore) {}

  readonly state$ = this.store.select(state => ({
    repository: state.repository,
    contributors: state.contributors.items,
    loading: state.contributors.fetching > 0,
  }));

  ngOnInit() {
    const repoFromUrl = new URLSearchParams(window.location.search).get('repo');
    if (repoFromUrl && repoFromUrl.trim().length > 0) {
      this.usecase.fetchContributors(repoFromUrl);
    } else {
      this.usecase.fetchContributors('angular/angular-ja');
    }
  }

  selectRepository(repoName: string) {
    this.usecase.fetchContributors(repoName);
  }
}
