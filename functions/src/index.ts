require('dotenv').config();

import * as cors from 'cors';
import * as admin from 'firebase-admin';
import * as functions from 'firebase-functions';

import { Repository } from './model/repository';
import { validateRepoParam } from './utils/validators';
import { fetchContributors } from './service/fetch-contributors';
import { renderContributorsImage } from './service/render-image';
import { ContributorsImageCache } from './service/cache-storage';
import { getApplicationConfig } from './service/app-config';

admin.initializeApp();

const bucket = admin.storage().bucket();
const config = getApplicationConfig();
console.debug('config', config);

export const createContributorsImage = functions
  .runWith({ timeoutSeconds: 60, memory: '1GB' })
  .https.onRequest(async (request, response) => {
    const repoParam = request.query['repo'];

    try {
      validateRepoParam(repoParam);
    } catch (error) {
      console.error(error);
      response.status(400).send(error.toString());
      return;
    }
    const repository = Repository.fromString(repoParam);
    console.debug(`repository: ${repository.toString()}`);

    const imageCache = new ContributorsImageCache(bucket, { useCache: config.useCache });

    const createImage = async () => {
      console.debug('restore cache');
      const cache = await imageCache.restore(repository);
      if (cache) {
        console.debug('cache hit');
        return cache;
      }
      console.debug(`render image`);
      return renderContributorsImage(repository, { webappUrl: config.webappUrl, useHeadless: config.useHeadless });
    };

    try {
      const image = await createImage();
      response
        .header('Content-Type', 'image/png')
        .header('Cache-Control', 'max-age=0, no-cache')
        .status(200)
        .send(image);
      console.debug('save cache');
      await imageCache.save(repository, image);
    } catch (error) {
      console.error(error);
      response.status(500).send(error);
    }
  });

const withCors = cors({ origin: true });
export const getContributors = functions.https.onRequest((request, response) => {
  withCors(request, response, async () => {
    const repoParam = request.query['repo'];

    try {
      validateRepoParam(repoParam);
    } catch (error) {
      console.error(error);
      response.status(400).send(error.toString());
      return;
    }
    const repository = Repository.fromString(repoParam);
    console.debug(`repository: ${repository.toString()}`);

    try {
      console.debug(`fetch contributors starting`);
      const contributors = await fetchContributors(repository);
      console.debug(`fetch contributors finished`);

      response
        .header('Content-Type', 'application/json')
        .status(200)
        .send(contributors);
    } catch (err) {
      console.error(err);
      response.status(500).send(err.toString());
    }
  });
});
