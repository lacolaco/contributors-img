import * as compression from 'compression';
import * as express from 'express';
import * as admin from 'firebase-admin';
import * as functions from 'firebase-functions';
import { createContributorsImageQuery } from './query/contributors-image.query';
import { createContributorsQuery } from './query/contributors.query';
import { getApplicationConfig } from './service/app-config';
import { ContributorsImageCache } from './service/image-cache';
import { ContributorsJsonCache } from './service/json-cache';
import { RepoInfoRepository } from './service/repo-info.repository';
import { Repository } from '@lib/core';
import { validateRepoParam } from './utils/validators';

admin.initializeApp();

// eslint-disable-next-line @typescript-eslint/no-explicit-any
const bucket = admin.storage().bucket() as any;
const config = getApplicationConfig();
console.debug('config', config);
const repoInfoRepository = new RepoInfoRepository(admin.firestore());
const contributorsJsonCache = new ContributorsJsonCache(bucket, { useCache: config.useCache });
const contributorsImageCache = new ContributorsImageCache(bucket, { useCache: config.useCache });

const contributorsQuery = createContributorsQuery(contributorsJsonCache);
const contributorsImageQuery = createContributorsImageQuery(contributorsImageCache, config);

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

    try {
      const contributors = await contributorsQuery(repository);
      const image = await contributorsImageQuery(repository, contributors);
      response
        .header('Content-Type', 'image/png')
        .header('Cache-Control', `max-age=${60 * 60}`)
        .status(200)
        .send(image);

      await repoInfoRepository.set(repository, {
        lastGeneratedAt: new Date(),
      });
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

      try {
        const contributors = await contributorsQuery(repository);
        response
          .header('Content-Type', 'application/json')
          .header('Cache-Control', `max-age=${60 * 60}`)
          .status(200)
          .send(contributors);
      } catch (err) {
        console.error(err);
        response.status(500).send(err.toString());
      }
    }),
);
