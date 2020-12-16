import { assertRepositoryName } from '@lib/core';
import { Request, Response } from 'express';
import { injectable } from 'tsyringe';
import { ContributorsQuery } from '../query/contributors';
import { Controller } from '../utils/types';
import { addTracingLabels, runWithTracing } from '../utils/tracing';

@injectable()
export class GetContributorsController implements Controller {
  constructor(private readonly contributorsQuery: ContributorsQuery) {}

  async onRequest(req: Request, res: Response) {
    const repoName = req.query.repo;
    if (!assertRepositoryName(repoName)) {
      res.status(400).send(`"${repoName}" is not a valid repository name`);
      return;
    }
    addTracingLabels({ 'app/repoName': repoName });
    try {
      const contributors = await runWithTracing('getContributors', () =>
        this.contributorsQuery.getContributors(repoName),
      );
      res.header('cache-control', `public, max-age=${60 * 60}`).json(contributors);
    } catch (err) {
      console.error(err);
      res.status(500).send(err.toString());
    }
  }
}
