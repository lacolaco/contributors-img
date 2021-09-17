import { assertRepositoryName, Repository } from '@lib/core';
import { Request, Response } from 'express';
import { injectable } from 'tsyringe';
import { GetContributorsUsecase } from '../usecase/contributors';
import { addTracingLabels } from '../utils/tracing';
import { Controller } from '../utils/types';

type Params = {
  repo: string;
  max?: string;
};

@injectable()
export class GetContributorsController implements Controller {
  constructor(private readonly usecase: GetContributorsUsecase) {}

  async onRequest(req: Request, res: Response) {
    const { repo, max } = req.query as Params;
    const maxOrNull = max ? Number(max) : null;
    // request validation
    if (!assertRepositoryName(repo)) {
      res.status(400).send(`"${repo}" is not a valid repository name`);
      return;
    }
    if (max != null && Number(max) > 0) {
      res.status(400).send(`max ${max} is not a positive integer`);
      return;
    }

    addTracingLabels({ 'app/repoName': repo });
    try {
      const contributors = await this.usecase.execute(Repository.fromString(repo), maxOrNull);
      res.header('cache-control', `public, max-age=${60 * 60}`).json(contributors);
    } catch (err) {
      console.error(err);
      if (err instanceof Error) {
        res.status(500).send(err.toString());
      } else {
        res.status(500).send('error');
      }
    }
  }
}
