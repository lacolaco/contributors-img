import { CommonModule } from '@angular/common';
import { Component, inject, OnDestroy } from '@angular/core';
import { Observable, Subject, takeUntil } from 'rxjs';
import { FeaturedRepositoryDatasourceToken } from '../../../../shared/featured-repository';
import { FeaturedRepository } from '../../../../shared/model/repository';
import { RepositoryGalleryComponent } from '../../component/repository-gallery/repository-gallery.component';

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
  private readonly datasource = inject(FeaturedRepositoryDatasourceToken);
  private readonly onDestroy$ = new Subject<void>();

  readonly repositories$: Observable<FeaturedRepository[]> = this.datasource.repositories$.pipe(
    takeUntil(this.onDestroy$),
  );

  ngOnDestroy() {
    this.onDestroy$.next();
  }
}
