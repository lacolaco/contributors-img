import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { ReactiveFormsModule } from '@angular/forms';
import { MatButtonModule } from '@angular/material/button';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatInputModule } from '@angular/material/input';
import { RxSubscribeModule } from '@soundng/rx-subscribe';
import { ContributorsListModule } from '../../shared/contributors-list/contributors-list.module';
import { FooterComponent } from './component/footer/footer.component';
import { HeaderComponent } from './component/header/header.component';
import { ImageSnippetComponent } from './component/image-snippet/image-snippet.component';
import { RepositoryFormComponent } from './component/repository-form/repository-form.component';
import { PreviewRoutingModule } from './preview-routing.module';
import { PreviewComponent } from './preview.component';

@NgModule({
  declarations: [PreviewComponent, HeaderComponent, FooterComponent, RepositoryFormComponent, ImageSnippetComponent],
  imports: [
    CommonModule,
    PreviewRoutingModule,
    ReactiveFormsModule,
    RxSubscribeModule,
    MatFormFieldModule,
    MatInputModule,
    MatButtonModule,
    ContributorsListModule,
  ],
})
export class PreviewModule {}
