import * as cors from 'cors';
import * as admin from 'firebase-admin';
import * as functions from 'firebase-functions';

import { Repository } from './model/repository';
import { validateRepoParam } from './utils/validators';
import { fetchContributors } from './service/fetch-contributors';
import { renderContributorsImage } from './service/render-image';
import { ContributorsImageCacheForCloudStorage, ContributorsImageCacheForLocal } from './service/cache-storage';

admin.initializeApp();

const bucket = admin.storage().bucket();

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
    console.log(`repository: ${repository.toString()}`);

    const imageCache =
      request.hostname !== 'localhost'
        ? new ContributorsImageCacheForCloudStorage(bucket)
        : new ContributorsImageCacheForLocal();

    const createImage = async () => {
      console.debug('restore cache');
      const cache = await imageCache.restore(repository);
      if (cache) {
        console.debug('cache hit');
        return cache;
      }
      console.log(`render image`);
      return renderContributorsImage(repoParam);
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
    console.log(`repository: ${repository.toString()}`);

    try {
      console.log(`fetch contributors starting`);
      const contributors = await fetchContributors(repository);
      console.log(`fetch contributors finished`);

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
