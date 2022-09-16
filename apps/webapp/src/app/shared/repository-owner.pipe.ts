import { Pipe, PipeTransform } from '@angular/core';
import { Repository } from './model/repository';

@Pipe({
  name: 'repositoryOwner',
})
export class RepositoryOwnerPipe implements PipeTransform {
  transform(value: string): string {
    return Repository.fromString(value).owner;
  }
}
