import { Pipe, PipeTransform } from '@angular/core';

@Pipe({
  name: 'avatarUrl',
})
export class AvatarUrlPipe implements PipeTransform {
  transform(original: string, size: number = 100): string {
    const url = new URL(original);
    url.searchParams.set('size', String(size));
    return url.toString();
  }
}
