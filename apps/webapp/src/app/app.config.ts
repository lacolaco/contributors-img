import { provideHttpClient, withXhr } from '@angular/common/http';
import { ApplicationConfig, provideZoneChangeDetection } from '@angular/core';
import { provideAnimations } from '@angular/platform-browser/animations';
import { provideRouter } from '@angular/router';
import { routes } from './app-routes';
import { provideFeaturedRepositoryDatasource } from './shared/featured-repository/firestore';
import { provideFirebase } from './shared/firebase';

export const appConfig: ApplicationConfig = {
  providers: [
    // Angular 21 defaults ZONELESS_ENABLED to true. This keeps the app zone-based; removing it is not a cleanup.
    provideZoneChangeDetection(),
    provideAnimations(),
    // Angular 22 made FetchBackend the default. This keeps the app on XHR; removing it is not a cleanup.
    provideHttpClient(withXhr()),
    provideRouter(routes),
    provideFirebase(),
    provideFeaturedRepositoryDatasource(),
  ],
};
