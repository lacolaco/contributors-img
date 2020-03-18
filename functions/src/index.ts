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
import { Contributor, Repository } from './shared/model';
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

    const send = (image: Buffer) => {
      response
        .header('Content-Type', 'image/png')
        .header('Cache-Control', 'max-age=0, no-cache')
        .status(200)
        .send(image);
    };

    try {
      console.debug('restore cache');
      const cached = await cache.restore(repository);
      if (cached) {
        console.debug('cache hit');
        send(cached);
        return;
      }
      console.debug(`render image`);
      const rendered = await renderContributorsImage(repository, {
        webappUrl: config.webappUrl,
        useHeadless: config.useHeadless,
      });
      send(rendered);
      console.debug('save cache');
      await cache.save(repository, rendered);
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

      const send = (data: Contributor[]) => {
        response
          .header('Content-Type', 'application/json')
          .status(200)
          .send(data);
      };

      try {
        console.debug('restore cache');
        const cached = await cache.restore(repository);
        if (cached) {
          console.debug('cache hit');
          send(cached);
          return;
        }
        console.debug(`fetch data`);

        const contributors = await fetchContributors(repository);
        send(contributors);
        console.debug('save cache');
        await cache.save(repository, contributors);
      } catch (err) {
        console.error(err);
        response.status(500).send(err.toString());
      }
    }),
);
