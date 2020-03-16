import { ChangeDetectionStrategy, Component, Input } from '@angular/core';
import { Contributor } from '@api/shared/model';

@Component({
  selector: 'app-contributor-list',
  templateUrl: './contributor-list.component.html',
  styleUrls: ['./contributor-list.component.scss'],
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class ContributorListComponent {
  @Input() list: Contributor[] = [];
}
