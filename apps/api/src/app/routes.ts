import { Router } from 'express';
import { container } from 'tsyringe';
import { GetContributorsController } from './controller/get-contributors';
import { GetImageController } from './controller/get-image';
import { GetImage2Controller } from './controller/get-image2';

export default (): Router => {
  const router = Router();

  const getContributors = container.resolve(GetContributorsController);
  const getImage = container.resolve(GetImageController);
  const getImage2 = container.resolve(GetImage2Controller);

  router.get('/api/contributors', (req, res) => getContributors.onRequest(req, res));
  router.get('/api/svg', (req, res) => getImage2.onRequest(req, res));
  router.get('/image', (req, res) => getImage.onRequest(req, res));

  return router;
};
