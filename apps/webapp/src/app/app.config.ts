import { provideHttpClient } from '@angular/common/http';
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
    provideHttpClient(),
    provideRouter(routes),
    provideFirebase(),
    provideFeaturedRepositoryDatasource(),
  ],
};
