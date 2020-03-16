import { HttpClientModule } from '@angular/common/http';
import { NgModule } from '@angular/core';
import { ReactiveFormsModule } from '@angular/forms';
import { MatButtonModule } from '@angular/material/button';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatInputModule } from '@angular/material/input';
import { BrowserModule } from '@angular/platform-browser';
import { BrowserAnimationsModule } from '@angular/platform-browser/animations';
import { RxSubscribeModule } from '@soundng/rx-subscribe';
import { AppComponent } from './app.component';
import { ContributorListComponent } from './view/component/contributor-list/contributor-list.component';
import { FooterComponent } from './view/component/footer/footer.component';
import { HeaderComponent } from './view/component/header/header.component';
import { ImageSnippetComponent } from './view/component/image-snippet/image-snippet.component';
import { RepositoryFormComponent } from './view/component/repository-form/repository-form.component';

@NgModule({
  declarations: [
    AppComponent,
    HeaderComponent,
    RepositoryFormComponent,
    ContributorListComponent,
    ImageSnippetComponent,
    FooterComponent,
  ],
  imports: [
    BrowserModule,
    HttpClientModule,
    ReactiveFormsModule,
    BrowserAnimationsModule,
    RxSubscribeModule,
    MatFormFieldModule,
    MatInputModule,
    MatButtonModule,
  ],
  providers: [],
  bootstrap: [AppComponent],
})
export class AppModule {}
