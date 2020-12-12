import { Firestore } from '@google-cloud/firestore';

export const createFirestore = (): Firestore => {
  return new Firestore({});
};
