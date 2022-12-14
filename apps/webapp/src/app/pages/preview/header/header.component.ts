import { ChangeDetectionStrategy, Component } from '@angular/core';
import { MatLegacyButtonModule as MatButtonModule } from '@angular/material/legacy-button';

@Component({
  selector: 'app-header',
  template: `
    <h1 class="heading">contrib.rocks</h1>
    <div class="discription">Generate an image of contributors to keep your README.md in sync.</div>
    <div class="banner">
      <a mat-stroked-button href="https://github.com/sponsors/lacolaco?frequency=one-time">Buy author a coffee</a>
    </div>
  `,
  styleUrls: ['./header.component.scss'],
  changeDetection: ChangeDetectionStrategy.OnPush,
  standalone: true,
  imports: [MatButtonModule],
})
export class HeaderComponent {}
