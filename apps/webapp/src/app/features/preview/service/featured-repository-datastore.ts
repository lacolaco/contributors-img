import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';
import { FeaturedRepository } from '../../../shared/model/repository';

@Injectable()
export abstract class FeaturedRepositoryDatastore {
  abstract readonly repositories$: Observable<FeaturedRepository[]>;
}
