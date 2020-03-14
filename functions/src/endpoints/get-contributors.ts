import * as cors from 'cors';
import * as functions from 'firebase-functions';
import { Repository } from '../model/repository';
import { fetchContributors } from '../service/fetch-contributors';
import { validateRepoParam } from './utils/validators';
import { createEndpoint } from './utils/create-endpoint';

const withCors = cors({ origin: true });

export const getContributors = createEndpoint(app =>
  functions.https.onRequest((request, response) => {
    withCors(request, response, async () => {
      console.group('getContributors');
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
        const contributors = await fetchContributors(repository);
        response
          .header('Content-Type', 'application/json')
          .status(200)
          .send(contributors);
      } catch (error) {
        console.error(error);
        response.status(500).send(error.toString());
      }
      console.groupEnd();
    });
  }),
);
