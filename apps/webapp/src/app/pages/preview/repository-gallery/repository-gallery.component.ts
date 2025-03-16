import { CommonModule } from '@angular/common';
import { ChangeDetectionStrategy, Component, input } from '@angular/core';
import { RouterLink } from '@angular/router';
import { FeaturedRepository } from '../../../models/repository';
import { RepositoryImageUrlPipe } from './repository-image-url.pipe';

@Component({
  selector: 'app-repository-gallery',
  templateUrl: './repository-gallery.component.html',
  styleUrls: ['./repository-gallery.component.scss'],
  changeDetection: ChangeDetectionStrategy.OnPush,
  imports: [CommonModule, RouterLink, RepositoryImageUrlPipe],
})
export class RepositoryGalleryComponent {
  readonly repositories = input<FeaturedRepository[]>([]);
}
