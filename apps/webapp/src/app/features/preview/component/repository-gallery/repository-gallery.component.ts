import { ChangeDetectionStrategy, Component, Input } from '@angular/core';

@Component({
  selector: 'app-repository-gallery',
  templateUrl: './repository-gallery.component.html',
  styleUrls: ['./repository-gallery.component.scss'],
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class RepositoryGalleryComponent {
  @Input()
  repositories: { name: string }[];
}
