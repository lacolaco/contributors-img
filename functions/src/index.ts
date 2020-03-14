import * as admin from 'firebase-admin';
import { HttpsFunction } from 'firebase-functions';

admin.initializeApp();

import * as endpoints from './endpoints';

export const getContributors: HttpsFunction = endpoints.getContributors(admin);
export const createContributorsImage: HttpsFunction = endpoints.createContributorsImage(admin);
