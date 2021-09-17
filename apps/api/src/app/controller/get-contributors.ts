import { assertRepositoryName } from '@lib/core';
import { Request, Response } from 'express';
import { injectable } from 'tsyringe';
import { GetContributorsUsecase } from '../usecase/contributors';
import { addTracingLabels } from '../utils/tracing';
import { Controller } from '../utils/types';

@injectable()
export class GetContributorsController implements Controller {
  constructor(private readonly getContributorsUsecase: GetContributorsUsecase) {}

  async onRequest(req: Request, res: Response) {
    const repoName = req.query['repo'];
    const maxCount = Number(req.query['max']) || null;
    if (!assertRepositoryName(repoName)) {
      res.status(400).send(`"${repoName}" is not a valid repository name`);
      return;
    }
    if (typeof maxCount === 'number' && maxCount < 1) {
      res.status(400).send('max must be a positive integer');
      return;
    }

    addTracingLabels({ 'app/repoName': repoName });
    try {
      const contributors = await this.getContributorsUsecase.execute(repoName, { maxCount: maxCount ?? 100 });
      res.header('cache-control', `public, max-age=${60 * 60}`).json(contributors);
    } catch (err) {
      console.error(err);
      res.status(500).send(err.toString());
    }
  }
}
