import { Component, OnInit } from '@angular/core';
import { AppStore } from './state/store';
import { FetchContributorsUsecase } from './usecase/fetch-contributors.usecase';
import { requestChangeDetection } from './utils/rx/detect-changes';

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.scss'],
})
export class AppComponent implements OnInit {
  constructor(private fetchContributorsUsecase: FetchContributorsUsecase, private store: AppStore) {}

  readonly state$ = this.store
    .select(state => ({
      repository: state.repository,
      contributors: state.contributors.items,
      loading: state.contributors.fetching > 0,
    }))
    .pipe(requestChangeDetection(this));

  ngOnInit() {
    const repoFromUrl = new URLSearchParams(window.location.search).get('repo');
    if (repoFromUrl && repoFromUrl.trim().length > 0) {
      this.fetchContributorsUsecase.execute(repoFromUrl);
    } else {
      this.fetchContributorsUsecase.execute('angular/angular-ja');
    }
  }

  selectRepository(repoName: string) {
    this.fetchContributorsUsecase.execute(repoName);
  }
}
