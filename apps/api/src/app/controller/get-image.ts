import { assertRepositoryName } from '@lib/core';
import { Request, Response } from 'express';
import { injectable } from 'tsyringe';
import { Controller } from '../utils/types';
import { ContributorsImageQuery } from '../query/contributors-image';
import { addTracingLabels, runWithTracing } from '../utils/tracing';

@injectable()
export class GetImageController implements Controller {
  constructor(private readonly imageQuery: ContributorsImageQuery) {}
  async onRequest(req: Request, res: Response) {
    const repoName = req.query.repo;
    if (!assertRepositoryName(repoName)) {
      res.status(400).send(`"${repoName}" is not a valid repository name`);
      return;
    }
    addTracingLabels({ 'app/repoName': repoName });
    try {
      const acceptWebp = req.accepts('webp') !== false;
      const { fileStream, contentType } = await runWithTracing('getImage', () =>
        this.imageQuery.getImage(repoName, { acceptWebp }),
      );
      res
        .header('Content-Type', contentType)
        .header('Vary', `Accept`)
        .header('Cache-Control', `public, max-age=${60 * 60 * 6}`);
      fileStream.pipe(res);
    } catch (err) {
      console.error(err);
      res.status(500).send(err.toString());
    }
  }
}
