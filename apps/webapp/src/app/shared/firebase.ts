import { InjectionToken, Provider } from '@angular/core';
import { FirebaseApp, initializeApp } from 'firebase/app';
import { getAnalytics } from 'firebase/analytics';
import { Firestore, getFirestore } from 'firebase/firestore';
import { getPerformance } from 'firebase/performance';
import { environment } from '../../environments/environment';

export const FirestoreToken = new InjectionToken<Firestore>('Firestore');

export function provideFirebase(): Provider[] {
  const app: FirebaseApp = initializeApp(environment.firebaseConfig);
  // Not consumed directly; initializing them registers auto-instrumentation.
  getAnalytics(app);
  getPerformance(app);

  return [{ provide: FirestoreToken, useValue: getFirestore(app) }];
}
