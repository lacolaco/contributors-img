import { CommonModule } from '@angular/common';
import { Component, OnDestroy } from '@angular/core';
import { Observable, Subject, takeUntil } from 'rxjs';
import { FeaturedRepository } from '../../../../shared/model/repository';
import { RepositoryGalleryComponent } from '../../component/repository-gallery/repository-gallery.component';
import { FeaturedRepositoryDatastore } from '../../service/featured-repository-datastore';

@Component({
  selector: 'app-recent-usage',
  template: `
    <ng-container *ngIf="repositories$ | async as repositories">
      <app-repository-gallery [repositories]="repositories"></app-repository-gallery>
    </ng-container>
  `,
  styles: [
    `
      :host {
        display: block;
        width: 100%;
      }
    `,
  ],
  standalone: true,
  imports: [CommonModule, RepositoryGalleryComponent],
})
export class RecentUsageComponent implements OnDestroy {
  private readonly onDestroy$ = new Subject<void>();

  readonly repositories$: Observable<FeaturedRepository[]> = this.featuredRepositories.repositories$.pipe(
    takeUntil(this.onDestroy$),
  );

  constructor(private readonly featuredRepositories: FeaturedRepositoryDatastore) {}

  ngOnDestroy() {
    this.onDestroy$.next();
  }
}
