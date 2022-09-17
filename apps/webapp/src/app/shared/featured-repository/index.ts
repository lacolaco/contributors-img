import { InjectionToken } from '@angular/core';
import { Observable } from 'rxjs';
import { FeaturedRepository } from '../../models';

export const FeaturedRepositoryDatasourceToken = new InjectionToken<FeaturedRepositoryDatasource>(
  'FeaturedRepositoryDatasource',
);

export abstract class FeaturedRepositoryDatasource {
  abstract readonly repositories$: Observable<FeaturedRepository[]>;
}
