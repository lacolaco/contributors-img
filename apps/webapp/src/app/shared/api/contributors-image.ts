import { HttpClient, HttpParams } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Repository } from '@lib/core';

@Injectable({ providedIn: 'root' })
export class ContributorsImageApi {
  constructor(private readonly httpClient: HttpClient) {}

  getByRepository(repository: Repository, { max, columns }: { max: number | null; columns: number | null }) {
    const params = new HttpParams({
      fromObject: {
        repo: repository.toString(),
        preview: true,
        max: max ?? '',
        columns: columns ?? '',
      },
    });
    return this.httpClient.get('/image', {
      responseType: 'text',
      params,
    });
  }
}
