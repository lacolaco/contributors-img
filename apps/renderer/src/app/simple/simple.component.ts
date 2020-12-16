import { ChangeDetectionStrategy, Component, Input } from '@angular/core';
import { Contributor } from '@lib/core';

@Component({
  selector: 'app-renderer-simple',
  templateUrl: './simple.component.html',
  styleUrls: ['./simple.component.scss'],
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class SimpleRendererComponent {
  @Input()
  contributors: Contributor[];
}
