import * as functions from 'firebase-functions';
import { Repository } from '../model/repository';
import { generateContributorsImage } from '../service/contributors-image';
import { restoreImageCache, saveImageCache } from '../service/image-cache';
import { validateRepoParam } from './utils/validators';

export const createContributorsImage = functions
  .runWith({ timeoutSeconds: 60, memory: '1GB' })
  .https.onRequest(async (request, response) => {
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

    try {
      const image = await createImage(repository);
      response
        .header('Content-Type', 'image/png')
        .header('Cache-Control', 'max-age=0, no-cache')
        .status(200)
        .send(image);
      console.debug('save cache');
      await saveImageCache(repository, image);
    } catch (error) {
      console.error(error);
      response.status(500).send(error.toString());
    }
    console.groupEnd();
  });

async function createImage(repository: Repository) {
  console.debug('restore cache');
  const cache = await restoreImageCache(repository);
  if (cache) {
    console.debug('cache is found');
    return cache;
  }
  console.debug('generate image');
  const image = await generateContributorsImage(repository);
  return image;
}
