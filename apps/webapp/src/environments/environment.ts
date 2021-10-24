import { firebaseConfig } from './firebase-config';

export const environment = {
  production: false,
  firebaseConfig: firebaseConfig.prod,
  firestoreRootCollectionName: 'production',
};
