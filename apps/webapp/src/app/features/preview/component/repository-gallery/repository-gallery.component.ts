import { ChangeDetectionStrategy, Component, Input } from '@angular/core';
import { FeaturedRepository } from '../../../../shared/model/repository';

@Component({
  selector: 'app-repository-gallery',
  templateUrl: './repository-gallery.component.html',
  styleUrls: ['./repository-gallery.component.scss'],
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class RepositoryGalleryComponent {
  @Input()
  repositories: FeaturedRepository[] = [];
}
