import { Provider } from '@angular/core';
import { of } from 'rxjs';
import { FeaturedRepositoryDatasource, FeaturedRepositoryDatasourceToken } from './index';

class NoopFeaturedRepositoryDatasource implements FeaturedRepositoryDatasource {
  readonly repositories$ = of([]);
}

export function provideNoopFeaturedRepositoryDatasource(): Provider[] {
  return [
    {
      provide: FeaturedRepositoryDatasourceToken,
      useFactory: () => new NoopFeaturedRepositoryDatasource(),
    },
  ];
}
