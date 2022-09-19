import { Pipe, PipeTransform } from '@angular/core';
import { FeaturedRepository, Repository } from '../../../models/repository';

@Pipe({
  name: 'repositoryImageUrl',
  standalone: true,
})
export class RepositoryImageUrlPipe implements PipeTransform {
  transform(value: FeaturedRepository): string {
    return `https://github.com/${Repository.fromString(value.repository).owner}.png?w=64`;
  }
}
