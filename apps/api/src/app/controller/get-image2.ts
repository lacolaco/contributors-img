import { assertRepositoryName } from '@lib/core';
import { Request, Response } from 'express';
import { injectable } from 'tsyringe';
import { ContributorsImageQuery } from '../query/contributors-image';
import { ContributorsImageRenderer } from '../service/contributors-image-renderer';
import { Controller } from '../utils/types';

@injectable()
export class GetImage2Controller implements Controller {
  constructor(
    private readonly imageQuery: ContributorsImageQuery,
    private readonly renderer: ContributorsImageRenderer,
  ) {}
  async onRequest(req: Request, res: Response) {
    const repoName = req.query.repo;
    if (!assertRepositoryName(repoName)) {
      res.status(400).send(`"${repoName}" is not a valid repository name`);
      return;
    }

    const svg = this.renderer.renderSvg([]);

    res.status(200).send(svg);
  }
}
