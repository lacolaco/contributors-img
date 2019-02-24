import { HttpClient } from '@angular/common/http';
import { Component } from '@angular/core';
import { Observable } from 'rxjs';

interface Contributor {
  avatar_url: string;
  contributions: number;
  html_url: string;
  id: number;
  login: string;
}

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.scss'],
})
export class AppComponent {
  contributors$: Observable<Contributor[]>;

  constructor(private http: HttpClient) {
    const repo = new URLSearchParams(window.location.search).get('repo');

    if (!repo) {
      throw new Error('repo parameter is required.');
    }

    this.contributors$ = this.http.get<Contributor[]>(
      `https://us-central1-contributors-img.cloudfunctions.net/getContributors`,
      { params: { repo } },
    );
  }
}
