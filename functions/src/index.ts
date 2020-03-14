import * as firebase from 'firebase-admin';
import { HttpsFunction } from 'firebase-functions';

const app = firebase.initializeApp();

import * as endpoints from './endpoints';

export const getContributors: HttpsFunction = endpoints.getContributors(app);
export const createContributorsImage: HttpsFunction = endpoints.createContributorsImage(app);
