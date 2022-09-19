import { CommonModule } from '@angular/common';
import { ChangeDetectionStrategy, Component, Input } from '@angular/core';
import { RouterLinkWithHref } from '@angular/router';
import { FeaturedRepository } from '../../../models/repository';
import { RepositoryImageUrlPipe } from './repository-image-url.pipe';

@Component({
  selector: 'app-repository-gallery',
  templateUrl: './repository-gallery.component.html',
  styleUrls: ['./repository-gallery.component.scss'],
  changeDetection: ChangeDetectionStrategy.OnPush,
  standalone: true,
  imports: [CommonModule, RouterLinkWithHref, RepositoryImageUrlPipe],
})
export class RepositoryGalleryComponent {
  @Input()
  repositories: FeaturedRepository[] = [];
}
