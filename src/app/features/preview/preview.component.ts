import { Component, OnInit } from '@angular/core';
import { PreviewStore } from './store';
import { FetchContributorsUsecase } from './usecase/fetch-contributors.usecase';

@Component({
  selector: 'app-preview',
  templateUrl: './preview.component.html',
  styleUrls: ['./preview.component.scss'],
})
export class PreviewComponent implements OnInit {
  constructor(private fetchContributors: FetchContributorsUsecase, private store: PreviewStore) {}

  readonly state$ = this.store.select(state => ({
    repository: state.repository,
    contributors: state.contributors.items,
    loading: state.contributors.fetching > 0,
  }));

  ngOnInit() {
    const repoFromUrl = new URLSearchParams(window.location.search).get('repo');
    if (repoFromUrl && repoFromUrl.trim().length > 0) {
      this.fetchContributors.execute(repoFromUrl);
    } else {
      this.fetchContributors.execute('angular/angular-ja');
    }
  }

  selectRepository(repoName: string) {
    this.fetchContributors.execute(repoName);
  }
}
