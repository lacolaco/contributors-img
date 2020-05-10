import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { AngularFireModule } from '@angular/fire';
import { AngularFirestoreModule } from '@angular/fire/firestore';
import { ReactiveFormsModule } from '@angular/forms';
import { MatButtonModule } from '@angular/material/button';
import { MatCardModule } from '@angular/material/card';
import { MatDialogModule } from '@angular/material/dialog';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatInputModule } from '@angular/material/input';
import { ReactiveComponentModule } from '@ngrx/component';
import { environment } from '../../../environments/environment';
import { ContributorsListModule } from '../../shared/contributors-list/contributors-list.module';
import { FooterComponent } from './component/footer/footer.component';
import { HeaderComponent } from './component/header/header.component';
import { ImageSnippetComponent } from './component/image-snippet/image-snippet.component';
import { RepositoryFormComponent } from './component/repository-form/repository-form.component';
import { RepositoryGalleryComponent } from './component/repository-gallery/repository-gallery.component';
import { PreviewRoutingModule } from './preview-routing.module';
import { PreviewComponent } from './preview.component';

import 'firebase/firestore';

@NgModule({
  declarations: [
    PreviewComponent,
    HeaderComponent,
    FooterComponent,
    RepositoryFormComponent,
    ImageSnippetComponent,
    RepositoryGalleryComponent,
  ],
  imports: [
    CommonModule,
    PreviewRoutingModule,
    AngularFireModule.initializeApp(environment.firebaseConfig),
    AngularFirestoreModule,
    ReactiveFormsModule,
    ReactiveComponentModule,
    MatFormFieldModule,
    MatInputModule,
    MatButtonModule,
    MatCardModule,
    MatDialogModule,
    ContributorsListModule,
  ],
})
export class PreviewModule {}
