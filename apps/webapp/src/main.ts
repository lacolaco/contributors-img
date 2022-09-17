import { HttpClientModule } from '@angular/common/http';
import { enableProdMode, importProvidersFrom } from '@angular/core';
import { bootstrapApplication } from '@angular/platform-browser';
import { provideAnimations } from '@angular/platform-browser/animations';
import { AppComponent } from './app/app.component';

import { AngularFireModule } from '@angular/fire/compat';
import { AngularFireAnalyticsModule } from '@angular/fire/compat/analytics';
import { AngularFirestoreModule } from '@angular/fire/compat/firestore';
import { AngularFirePerformanceModule } from '@angular/fire/compat/performance';
import { provideRouter } from '@angular/router';
import { routes } from './app/app-routing';
import { FeaturedRepositoryDatastore } from './app/features/preview/service/featured-repository-datastore';
import { FeaturedRepositoryDatastoreImpl } from './app/features/preview/service/featured-repository-datastore.impl';
import { environment } from './environments/environment';

if (environment.production) {
  enableProdMode();
}

bootstrapApplication(AppComponent, {
  providers: [
    provideAnimations(),
    importProvidersFrom(HttpClientModule),
    provideRouter(routes),
    importProvidersFrom(AngularFireModule.initializeApp(environment.firebaseConfig)),
    importProvidersFrom(AngularFirestoreModule),
    importProvidersFrom(AngularFireAnalyticsModule),
    importProvidersFrom(AngularFirePerformanceModule),
    {
      provide: FeaturedRepositoryDatastore,
      useClass: FeaturedRepositoryDatastoreImpl,
    },
  ],
}).catch((err) => console.error(err));
