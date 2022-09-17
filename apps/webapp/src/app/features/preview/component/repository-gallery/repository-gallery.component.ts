import { CommonModule } from '@angular/common';
import { ChangeDetectionStrategy, Component, Input } from '@angular/core';
import { RouterLinkWithHref } from '@angular/router';
import { RepositoryOwnerPipe } from '../../../../shared/repository-owner.pipe';
import { FeaturedRepository } from '../../../../shared/model/repository';

@Component({
  selector: 'app-repository-gallery',
  templateUrl: './repository-gallery.component.html',
  styleUrls: ['./repository-gallery.component.scss'],
  changeDetection: ChangeDetectionStrategy.OnPush,
  standalone: true,
  imports: [CommonModule, RouterLinkWithHref, RepositoryOwnerPipe],
})
export class RepositoryGalleryComponent {
  @Input()
  repositories: FeaturedRepository[] = [];
}
