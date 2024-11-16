import { ChangeDetectionStrategy, Component } from '@angular/core';

@Component({
  selector: 'app-footer',
  template: `
    <div class="container">
      <span>
        &copy; 2019 - Suguru Inatomi &#64;
        <a target="_blank" rel="noopener" href="https://twitter.com/laco2net">lacolaco</a>
      </span>
      <div>
        <a target="_blank" rel="noopener" href="https://github.com/lacolaco/contributors-img">GitHub</a> /
        <a target="_blank" rel="noopener" href="https://github.com/sponsors/lacolaco">Become a sponsor</a>
      </div>
    </div>
  `,
  styleUrls: ['./footer.component.scss'],
  changeDetection: ChangeDetectionStrategy.OnPush,
  standalone: true,
})
export class FooterComponent {}
