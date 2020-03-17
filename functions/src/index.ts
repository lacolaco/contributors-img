require('dotenv').config();

import * as compression from 'compression';
import * as express from 'express';
import * as admin from 'firebase-admin';
import * as functions from 'firebase-functions';
import { getApplicationConfig } from './service/app-config';
import { fetchContributors } from './service/fetch-contributors';
import { ContributorsImageCache } from './service/image-cache';
import { ContributorsJsonCache } from './service/json-cache';
import { renderContributorsImage } from './service/render-image';
import { Repository } from './shared/model/repository';
import { validateRepoParam } from './utils/validators';

admin.initializeApp();

const bucket = admin.storage().bucket();
const config = getApplicationConfig();
console.debug('config', config);

export const createContributorsImage = functions.runWith({ timeoutSeconds: 60, memory: '1GB' }).https.onRequest(
  express().get('*', async (request, response) => {
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

    const cache = new ContributorsImageCache(bucket, { useCache: config.useCache });

    const prepare = async () => {
      console.debug('restore cache');
      const cached = await cache.restore(repository);
      if (cached) {
        console.debug('cache hit');
        return cached;
      }
      console.debug(`render image`);
      return renderContributorsImage(repository, { webappUrl: config.webappUrl, useHeadless: config.useHeadless });
    };

    try {
      const image = await prepare();
      response
        .header('Content-Type', 'image/png')
        .header('Cache-Control', 'max-age=0, no-cache')
        .status(200)
        .send(image);
      console.debug('save cache');
      await cache.save(repository, image);
    } catch (error) {
      console.error(error);
      response.status(500).send(error);
    }
  }),
);

export const getContributors = functions.https.onRequest(
  express()
    .use(compression({}))
    .get('*', async (request, response) => {
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
      const cache = new ContributorsJsonCache(bucket, { useCache: config.useCache });

      const prepare = async () => {
        console.debug('restore cache');
        const cached = await cache.restore(repository);
        if (cached) {
          console.debug('cache hit');
          return cached;
        }
        console.debug(`fetch data`);
        return fetchContributors(repository);
      };

      try {
        const contributors = await prepare();
        response
          .header('Content-Type', 'application/json')
          .status(200)
          .send(contributors);
        console.debug('save cache');
        await cache.save(repository, contributors);
      } catch (err) {
        console.error(err);
        response.status(500).send(err.toString());
      }
    }),
);
