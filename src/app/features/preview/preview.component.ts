import { Component, OnDestroy, OnInit } from '@angular/core';
import { AngularFirestore } from '@angular/fire/firestore';
import { ActivatedRoute } from '@angular/router';
import { BehaviorSubject, combineLatest, Subject } from 'rxjs';
import { map, takeUntil, throttleTime } from 'rxjs/operators';
import { PreviewStore } from './store';
import { FetchContributorsUsecase } from './usecase/fetch-contributors.usecase';

@Component({
  selector: 'app-preview',
  templateUrl: './preview.component.html',
  styleUrls: ['./preview.component.scss'],
})
export class PreviewComponent implements OnInit, OnDestroy {
  constructor(
    private route: ActivatedRoute,
    private fetchContributors: FetchContributorsUsecase,
    private store: PreviewStore,
    private firestore: AngularFirestore,
  ) {}

  private readonly showImageSnippetSubject = new BehaviorSubject<boolean>(false);
  private readonly onDestroy$ = new Subject();

  readonly state$ = combineLatest([
    this.store.valueChanges,
    this.firestore
      .collection<{ name: string }>('repositories', q => q.limit(12).orderBy('lastGeneratedAt', 'desc'))
      .valueChanges()
      .pipe(throttleTime(1000 * 10)),
    this.showImageSnippetSubject.asObservable(),
  ]).pipe(
    map(([state, repositories, showImageSnippet]) => ({
      repository: state.repository,
      contributors: state.contributors.items,
      loading: state.contributors.fetching > 0,
      repositories,
      showImageSnippet,
    })),
  );

  ngOnInit() {
    this.route.queryParamMap
      .pipe(
        takeUntil(this.onDestroy$),
        map(q => q.get('repo')),
      )
      .subscribe(repo => this.selectRepository(repo || 'angular/angular-ja'));
  }

  ngOnDestroy() {
    this.onDestroy$.next();
  }

  selectRepository(repoName: string) {
    this.showImageSnippetSubject.next(false);
    this.fetchContributors.execute(repoName);
  }

  showImageSnippet() {
    this.showImageSnippetSubject.next(true);
  }
}
