import { HttpClient } from '@angular/common/http';
import { Component, OnInit, ChangeDetectionStrategy } from '@angular/core';
import { DomSanitizer } from '@angular/platform-browser';
import { BehaviorSubject } from 'rxjs';

@Component({
  selector: 'app-preview2',
  templateUrl: './preview2.component.html',
  styleUrls: ['./preview2.component.scss'],
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class Preview2Component implements OnInit {
  svg$ = new BehaviorSubject(this.sanitizer.bypassSecurityTrustHtml(''));

  constructor(private readonly httpClient: HttpClient, private sanitizer: DomSanitizer) {}

  ngOnInit(): void {
    this.httpClient
      .get('/api/svg?repo=angular/angular-ja', { responseType: 'text' })
      .toPromise()
      .then((resp) => {
        this.svg$.next(this.sanitizer.bypassSecurityTrustHtml(resp));
      });
  }
}
