import { Router } from 'express';
import { container } from 'tsyringe';
import { GetContributorsController } from './controller/get-contributors';
import { GetImageController } from './controller/get-image';

export default (): Router => {
  const router = Router();

  const getContributors = container.resolve(GetContributorsController);
  const getImage = container.resolve(GetImageController);

  router.get('/api/contributors', (req, res) => getContributors.onRequest(req, res));
  router.get('/image', (req, res) => getImage.onRequest(req, res));

  return router;
};
