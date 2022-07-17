import { Injectable } from '@angular/core';
import { FeaturedRepository } from '@lib/core';
import { Observable } from 'rxjs';

@Injectable()
export abstract class FeaturedRepositoryDatastore {
  abstract readonly repositories$: Observable<FeaturedRepository[]>;
}
