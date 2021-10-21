import { Router } from 'express';
import { container } from 'tsyringe';
import { UpdateFeaturedRepositoriesController } from './controller/update-featured-repositories';

export default (): Router => {
  const router = Router();

  // TODO: request injector should be created per request
  const requestInjector = container.createChildContainer();

  const updateFeaturedRepositories = requestInjector.resolve(UpdateFeaturedRepositoriesController);
  router.get('/update-featured-repositories', (req, res) => updateFeaturedRepositories.onRequest(req, res));

  return router;
};
