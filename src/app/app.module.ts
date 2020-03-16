import { HttpClientModule } from '@angular/common/http';
import { NgModule } from '@angular/core';
import { ReactiveFormsModule } from '@angular/forms';
import { MatButtonModule } from '@angular/material/button';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatInputModule } from '@angular/material/input';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { BrowserModule } from '@angular/platform-browser';
import { RxSubscribeModule } from '@soundng/rx-subscribe';

import { AppComponent } from './app.component';
import { BrowserAnimationsModule } from '@angular/platform-browser/animations';
import { HeaderComponent } from './header/header.component';
import { FormComponent } from './form/form.component';
import { ContributorListComponent } from './contributor-list/contributor-list.component';
import { ImageSnippetComponent } from './view/component/image-snippet/image-snippet.component';
import { FooterComponent } from './view/component/footer/footer.component';

@NgModule({
  declarations: [AppComponent, HeaderComponent, FormComponent, ContributorListComponent, ImageSnippetComponent, FooterComponent],
  imports: [
    BrowserModule,
    HttpClientModule,
    ReactiveFormsModule,
    BrowserAnimationsModule,
    RxSubscribeModule,
    MatFormFieldModule,
    MatInputModule,
    MatButtonModule,
    MatProgressSpinnerModule,
  ],
  providers: [],
  bootstrap: [AppComponent],
})
export class AppModule {}
