import { provideHttpClient } from '@angular/common/http';
import { ApplicationConfig, provideZoneChangeDetection } from '@angular/core';
import { provideAnimations } from '@angular/platform-browser/animations';
import { provideRouter } from '@angular/router';
import { routes } from './app-routes';
import { provideFeaturedRepositoryDatasource } from './shared/featured-repository/firestore';
import { provideFirebase } from './shared/firebase';

export const appConfig: ApplicationConfig = {
  providers: [
    // Angular 21 defaults ZONELESS_ENABLED to true. This app drives change detection
    // through NgZone.run() in shared/featured-repository/firestore.ts, so it opts back
    // in explicitly. Dropping this switches the change-detection mode with no error.
    provideZoneChangeDetection(),
    provideAnimations(),
    provideHttpClient(),
    provideRouter(routes),
    provideFirebase(),
    provideFeaturedRepositoryDatasource(),
  ],
};
