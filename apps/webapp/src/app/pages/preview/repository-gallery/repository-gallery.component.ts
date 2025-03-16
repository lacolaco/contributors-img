import { CommonModule } from '@angular/common';
import { ChangeDetectionStrategy, Component, input } from '@angular/core';
import { RouterLink } from '@angular/router';
import { FeaturedRepository } from '../../../models/repository';
import { RepositoryImageUrlPipe } from './repository-image-url.pipe';

@Component({
  selector: 'app-repository-gallery',
  template: `
    <div class="heading">Used by</div>
    <ul class="gallery">
      @for (repo of repositories(); track repo) {
        <li class="gallery-item">
          <div class="repo">
            <img class="repo-image" height="32" width="32" src="{{ repo | repositoryImageUrl }}" />
            <div class="repo-body">
              <a href="https://github.com/{{ repo.repository }}" target="_blank" rel="noopener">
                <div class="repo-name">
                  <img style="display: flex" height="20" width="20" src="assets/images/github-64px.png" />
                  <span>{{ repo.repository }}</span>
                </div>
              </a>
              <div class="repo-numbers">
                <div>{{ repo.stargazers | number }} stars</div>
                <a style="font-weight: bold" routerLink [queryParams]="{ repo: repo.repository }">View rocks</a>
              </div>
            </div>
          </div>
        </li>
      }
    </ul>
  `,
  styleUrls: ['./repository-gallery.component.scss'],
  changeDetection: ChangeDetectionStrategy.OnPush,
  imports: [CommonModule, RouterLink, RepositoryImageUrlPipe],
})
export class RepositoryGalleryComponent {
  readonly repositories = input<FeaturedRepository[]>([]);
}
