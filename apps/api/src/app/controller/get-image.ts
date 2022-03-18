import { assertRepositoryName, Repository } from '@lib/core';
import { Request, Response } from 'express';
import { injectable } from 'tsyringe';
import { GetContributorsImageUsecase } from '../usecase/contributors-image';
import { isGitHubRequest } from '../utils/request';
import { addTracingLabels } from '../utils/tracing';
import { Controller } from '../utils/types';

type Params = {
  repo: string;
  max?: string;
  preview?: string;
  columns?: string;
};

@injectable()
export class GetImageController implements Controller {
  constructor(private readonly usecase: GetContributorsImageUsecase) {}
  async onRequest(req: Request, res: Response) {
    const { repo, max, preview, columns } = req.query as Params;

    const maxOrNull = max ? Number.parseInt(max, 10) : null;
    const columnsOrNull = columns ? Number.parseInt(columns, 10) : null;
    // request validation
    if (!assertRepositoryName(repo)) {
      res.status(400).send(`"${repo}" is not a valid repository name`);
      return;
    }
    if (maxOrNull != null && (!Number.isInteger(maxOrNull) || Number(max) < 1)) {
      res.status(400).send(`max ${max} is not a positive integer`);
      return;
    }
    if (columnsOrNull != null && (!Number.isInteger(columnsOrNull) || columnsOrNull < 1)) {
      res.status(400).send(`columns ${max} is not a positive integer`);
      return;
    }

    addTracingLabels({ 'app/repoName': repo });
    try {
      const fileStream = await this.usecase.execute({
        repository: Repository.fromString(repo),
        isGitHubRequest: isGitHubRequest(req),
        preview: !!preview,
        maxCount: maxOrNull,
        maxColumns: columnsOrNull,
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
