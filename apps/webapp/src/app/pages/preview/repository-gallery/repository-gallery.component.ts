import { CommonModule } from '@angular/common';
import { ChangeDetectionStrategy, Component, input } from '@angular/core';
import { RouterLink } from '@angular/router';
import { FeaturedRepository, Repository } from '../../../models/repository';

@Component({
  selector: 'app-repository-gallery',
  template: `
    <div class="mb-4 text-xl font-bold text-left w-full">Used by</div>
    <ul class="grid grid-cols-1 gap-4 md:grid-cols-2 lg:grid-cols-3 list-none w-full pl-0">
      @for (repo of repositories(); track repo) {
        <li class="box-border h-32">
          <div class="box-border w-full h-full p-4 border border-gray-300 rounded-lg flex flex-row">
            <img
              class="w-auto h-full aspect-square mr-4"
              height="32"
              width="32"
              [src]="getRepositoryImageUrl(repo)"
              alt="Repository contributors image for {{ repo.repository }}"
            />
            <div class="flex-1 flex flex-col justify-center">
              <a
                [href]="getRepositoryPageUrl(repo)"
                target="_blank"
                rel="noopener"
                class="text-black no-underline hover:underline"
              >
                <div class="flex items-center font-bold mb-4 text-base">
                  <img class="flex" height="20" width="20" src="assets/images/github-64px.png" alt="GitHub logo" />
                  <span class="ml-1">{{ repo.repository }}</span>
                </div>
              </a>
              <div class="flex justify-between">
                <div>{{ repo.stargazers | number }} stars</div>
                <a
                  class="font-bold text-black no-underline hover:underline"
                  routerLink
                  [queryParams]="{ repo: repo.repository }"
                  >View rocks</a
                >
              </div>
            </div>
          </div>
        </li>
      }
    </ul>
  `,
  changeDetection: ChangeDetectionStrategy.OnPush,
  imports: [CommonModule, RouterLink],
})
export class RepositoryGalleryComponent {
  readonly repositories = input<FeaturedRepository[]>([]);

  getRepositoryImageUrl(repo: FeaturedRepository): string {
    return `https://github.com/${Repository.fromString(repo.repository).owner}.png?w=64`;
  }

  getRepositoryPageUrl(repo: FeaturedRepository): string {
    return `https://github.com/${repo.repository}`;
  }
}
