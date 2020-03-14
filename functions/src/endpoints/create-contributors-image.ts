import * as functions from 'firebase-functions';
import { Repository } from '../model/repository';
import { generateContributorsImage } from '../service/contributors-image';
import { readFile, writeFile } from '../service/cache-storage';
import { validateRepoParam } from './utils/validators';
import { createEndpoint } from './utils/create-endpoint';

export const createContributorsImage = createEndpoint(app => {
  const cacheBucket = app.storage().bucket();
  return functions.runWith({ timeoutSeconds: 60, memory: '1GB' }).https.onRequest(async (request, response) => {
    console.group('createContributorsImage');
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
    const cacheId = generateCacheId(repository);

    const createImage = async () => {
      console.debug('restore cache');
      const cache = await readFile(cacheBucket, cacheId);
      if (cache) {
        console.debug('cache is found');
        return cache;
      }
      console.debug('generate image');
      const image = await generateContributorsImage(repository);
      return image;
    };

    try {
      const image = await createImage();
      response
        .header('Content-Type', 'image/png')
        .header('Cache-Control', 'max-age=0, no-cache')
        .status(200)
        .send(image);
      console.debug('save cache');
      await writeFile(cacheBucket, cacheId, image);
    } catch (error) {
      console.error(error);
      response.status(500).send(error.toString());
    }
    console.groupEnd();
  });
});

function generateCacheId(repository: Repository) {
  return `image-cache--${repository.owner}--${repository.repo}`;
}
