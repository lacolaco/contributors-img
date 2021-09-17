import { HttpClient, HttpParams } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Repository } from '@lib/core';

@Injectable({ providedIn: 'root' })
export class ContributorsImageApi {
  constructor(private readonly httpClient: HttpClient) {}

  getByRepository(repository: Repository, { max }: { max: number | null }) {
    const params = new HttpParams({
      fromObject: {
        repo: repository.toString(),
        preview: true,
        max: max ?? '',
      },
    });
    return this.httpClient.get('/image', {
      responseType: 'text',
      params,
    });
  }
}
