import { Component, Inject } from '@angular/core';
import { Contributor } from '@lib/core';
import { CONTRIBUTORS_DATA } from './tokens';

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.scss'],
})
export class AppComponent {
  constructor(@Inject(CONTRIBUTORS_DATA) public readonly contributors: Contributor[]) {
    if (!this.contributors) {
      throw new Error('No Embedded Contributors Data');
    }
  }
}
