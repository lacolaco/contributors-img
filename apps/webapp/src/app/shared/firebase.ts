import { InjectionToken, Provider } from '@angular/core';
import { FirebaseApp, initializeApp } from 'firebase/app';
import { Firestore, getFirestore } from 'firebase/firestore';
import { environment } from '../../environments/environment';

export const FirestoreToken = new InjectionToken<Firestore>('Firestore');

// Analytics and Performance are deliberately not started here. `provideAnalytics` /
// `providePerformance` only registered lazy factories, and nothing ever injected those
// tokens, so neither service was running before this migration. Starting them would be a
// behaviour change — and an unguarded one, since every environment file resolves to the
// production Firebase config.
export function provideFirebase(): Provider[] {
  const app: FirebaseApp = initializeApp(environment.firebaseConfig);

  return [{ provide: FirestoreToken, useValue: getFirestore(app) }];
}
