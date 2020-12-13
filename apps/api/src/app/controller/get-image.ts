import { assertRepositoryName } from '@lib/core';
import { Request, Response } from 'express';
import { injectable } from 'tsyringe';
import { Controller } from '../utils/types';
import { ContributorsImageQuery } from '../query/contributors-image';
import { runWithTracing } from '../utils/tracing';

@injectable()
export class GetImageController implements Controller {
  constructor(private readonly imageQuery: ContributorsImageQuery) {}
  async onRequest(req: Request, res: Response) {
    const repoName = req.query.repo;
    if (!assertRepositoryName(repoName)) {
      res.status(400).send(`"${repoName}" is not a valid repository name`);
      return;
    }
    try {
      const image = await runWithTracing('getImage', () => this.imageQuery.getImage(repoName));
      res
        .header('Content-Type', 'image/png')
        .header('Cache-Control', `max-age=${60 * 60}`)
        .send(image);
    } catch (err) {
      console.error(err);
      res.status(500).send(err.toString());
    }
  }
}
