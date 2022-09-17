import { ChangeDetectionStrategy, Component } from '@angular/core';

@Component({
  selector: 'app-header',
  template: `
    <div class="header">
      <div class="heading">contrib.rocks</div>
      <div class="discription">
        To Keep in sync your contributors list without any pain,<br />
        <b>ENTER</b> repository name and <b>GENERATE</b> an dynamic image URL for displaying it!
      </div>
    </div>
  `,
  styleUrls: ['./header.component.scss'],
  changeDetection: ChangeDetectionStrategy.OnPush,
  standalone: true,
})
export class HeaderComponent {}
