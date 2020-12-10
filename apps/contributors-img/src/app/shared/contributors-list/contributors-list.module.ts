import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { AvatarUrlPipe } from './pipe/avatar-url.pipe';
import { ContributorListComponent } from './component/contributor-list/contributor-list.component';

@NgModule({
  declarations: [AvatarUrlPipe, ContributorListComponent],
  imports: [CommonModule],
  exports: [AvatarUrlPipe, ContributorListComponent],
})
export class ContributorsListModule {}
