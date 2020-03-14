import * as firebase from 'firebase-admin';

firebase.initializeApp();

import * as endpoints from './endpoints';

export const getContributors = endpoints.getContributors;
export const createContributorsImage = endpoints.createContributorsImage;
