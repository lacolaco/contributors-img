import { NgModule } from '@angular/core';
import { BrowserModule } from '@angular/platform-browser';
import { ContributorsListModule } from '@lib/renderer-ui';
import { AppComponent } from './app.component';
import { SimpleRendererComponent } from './simple/simple.component';

@NgModule({
  declarations: [AppComponent, SimpleRendererComponent],
  imports: [BrowserModule.withServerTransition({ appId: 'serverApp' }), ContributorsListModule],
  providers: [],
  bootstrap: [AppComponent],
})
export class AppModule {}
