import { Router } from 'express';
import { container } from 'tsyringe';
import { GetImageController } from './controller/get-image';

export default (): Router => {
  const router = Router();

  // TODO: request injector should be created per request
  const requestInjector = container.createChildContainer();

  const getImage = requestInjector.resolve(GetImageController);
  router.get('/image', (req, res) => getImage.onRequest(req, res));

  return router;
};
