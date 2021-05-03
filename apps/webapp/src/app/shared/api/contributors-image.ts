import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Repository } from '@lib/core';

@Injectable({ providedIn: 'root' })
export class ContributorsImageApi {
  constructor(private readonly httpClient: HttpClient) {}

  getByRepository(repository: Repository) {
    return this.httpClient.get('/image', {
      params: {
        repo: repository.toString(),
        preview: true,
      },
      responseType: 'text',
    });
  }
}
