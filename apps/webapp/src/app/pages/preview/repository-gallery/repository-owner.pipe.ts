import { Pipe, PipeTransform } from '@angular/core';
import { Repository } from '../../../models/repository';

@Pipe({
  name: 'repositoryOwner',
  standalone: true,
})
export class RepositoryOwnerPipe implements PipeTransform {
  transform(value: string): string {
    return Repository.fromString(value).owner;
  }
}
