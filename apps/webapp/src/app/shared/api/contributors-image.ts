import { HttpClient, HttpParams } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { ImageParams } from '../../models/image-params';

@Injectable({ providedIn: 'root' })
export class ContributorsImageApi {
  constructor(private readonly httpClient: HttpClient) {}

  getImage({ repository, max, columns }: ImageParams) {
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
