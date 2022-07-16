import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { AngularFireModule } from '@angular/fire/compat';
import { AngularFireAnalyticsModule } from '@angular/fire/compat/analytics';
import { AngularFirestoreModule } from '@angular/fire/compat/firestore';
import { AngularFirePerformanceModule } from '@angular/fire/compat/performance';
import { ReactiveFormsModule } from '@angular/forms';
import { MatButtonModule } from '@angular/material/button';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatInputModule } from '@angular/material/input';
import { environment } from '../../../environments/environment';
import { SharedModule } from '../../shared/shared.module';
import { SvgViewModule } from '../../shared/svg-view/svg-view.component';
import { FooterComponent } from './component/footer/footer.component';
import { HeaderComponent } from './component/header/header.component';
import { ImageSnippetComponent } from './component/image-snippet/image-snippet.component';
import { RepositoryFormComponent } from './component/repository-form/repository-form.component';
import { RepositoryGalleryComponent } from './component/repository-gallery/repository-gallery.component';
import { RecentUsageComponent } from './container/recent-usage/recent-usage.component';
import { PreviewRoutingModule } from './preview-routing.module';
import { PreviewComponent } from './preview.component';
import { FeaturedRepositoryDatastore } from './service/featured-repository-datastore';
import { FeaturedRepositoryDatastoreImpl } from './service/featured-repository-datastore.impl';

@NgModule({
  declarations: [
    PreviewComponent,
    HeaderComponent,
    FooterComponent,
    RepositoryFormComponent,
    ImageSnippetComponent,
    RepositoryGalleryComponent,
    RecentUsageComponent,
  ],
  imports: [
    CommonModule,
    PreviewRoutingModule,
    AngularFireModule.initializeApp(environment.firebaseConfig),
    AngularFirestoreModule,
    AngularFireAnalyticsModule,
    AngularFirePerformanceModule,
    ReactiveFormsModule,
    MatFormFieldModule,
    MatInputModule,
    MatButtonModule,
    SvgViewModule,
    SharedModule,
  ],
  providers: [
    {
      provide: FeaturedRepositoryDatastore,
      useClass: FeaturedRepositoryDatastoreImpl,
    },
  ],
})
export class PreviewModule {}
