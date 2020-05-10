require('dotenv').config();

import * as express from 'express';
import * as compression from 'compression';
import { Storage } from '@google-cloud/storage';
import { Firestore } from '@google-cloud/firestore';

import { Repository } from '@contributors-img/api-interfaces';
import { validateRepoParam } from './app/utils/validators';
import { createContributorsQuery } from './app/query/contributors.query';
import { ContributorsJsonCache } from './app/service/json-cache';
import { getApplicationConfig } from './app/service/app-config';
import { RepoInfoRepository } from './app/service/repo-info.repository';
import { ContributorsImageCache } from './app/service/image-cache';
import { createContributorsImageQuery } from './app/query/contributors-image.query';

const config = getApplicationConfig();

const storage = new Storage();
const firestore = new Firestore();
const bucket = storage.bucket(`${storage.projectId}.appspot.com`);

console.debug('config', config);
const repoInfoRepository = new RepoInfoRepository(firestore);
const contributorsJsonCache = new ContributorsJsonCache(bucket, { useCache: config.useCache });
const contributorsImageCache = new ContributorsImageCache(bucket, { useCache: config.useCache });

const contributorsQuery = createContributorsQuery(contributorsJsonCache);
const contributorsImageQuery = createContributorsImageQuery(contributorsImageCache, config);

const app = express();
app.use(compression({}));

app.get('/image', async (request, response) => {
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
});

app.get('/api/contributors', async (request, response) => {
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
});

const port = process.env.PORT || 8080;
app.listen(port, () => {
  console.log('app listening on port', port);
});
