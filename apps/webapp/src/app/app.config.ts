import { provideHttpClient } from '@angular/common/http';
import { ApplicationConfig } from '@angular/core';
import { getAnalytics, provideAnalytics } from '@angular/fire/analytics';
import { initializeApp, provideFirebaseApp } from '@angular/fire/app';
import { getFirestore, provideFirestore } from '@angular/fire/firestore';
import { getPerformance, providePerformance } from '@angular/fire/performance';
import { provideAnimations } from '@angular/platform-browser/animations';
import { provideRouter } from '@angular/router';
import { environment } from '../environments/environment';
import { routes } from './app-routes';
import { provideFeaturedRepositoryDatasource } from './shared/featured-repository/firestore';

export const appConfig: ApplicationConfig = {
  providers: [
    provideAnimations(),
    provideHttpClient(),
    provideRouter(routes),
    provideFirebaseApp(() => initializeApp(environment.firebaseConfig)),
    provideFirestore(() => getFirestore()),
    provideAnalytics(() => getAnalytics()),
    providePerformance(() => getPerformance()),
    provideFeaturedRepositoryDatasource(),
  ],
};
