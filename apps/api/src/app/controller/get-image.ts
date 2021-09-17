import { assertRepositoryName, Repository } from '@lib/core';
import { Request, Response } from 'express';
import { injectable } from 'tsyringe';
import { GetContributorsImageUsecase } from '../usecase/contributors-image';
import { addTracingLabels } from '../utils/tracing';
import { Controller } from '../utils/types';

type Params = {
  repo: string;
  max?: string;
  preview?: string;
};

@injectable()
export class GetImageController implements Controller {
  constructor(private readonly usecase: GetContributorsImageUsecase) {}
  async onRequest(req: Request, res: Response) {
    const { repo, max, preview } = req.query as Params;

    const maxOrNull = max ? Number.parseInt(max, 10) : null;
    console.debug(maxOrNull);
    // request validation
    if (!assertRepositoryName(repo)) {
      res.status(400).send(`"${repo}" is not a valid repository name`);
      return;
    }
    if (maxOrNull != null && (!Number.isInteger(maxOrNull) || Number(max) < 1)) {
      res.status(400).send(`max ${max} is not a positive integer`);
      return;
    }

    addTracingLabels({ 'app/repoName': repo });
    try {
      const fileStream = await this.usecase.execute(Repository.fromString(repo), {
        preview: !!preview,
        maxCount: maxOrNull,
      });
      res
        .header('Content-Type', fileStream.contentType)
        .header('Vary', `Accept`)
        .header('Cache-Control', `public, max-age=${60 * 60 * 6}`);
      fileStream.data.pipe(res);
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
