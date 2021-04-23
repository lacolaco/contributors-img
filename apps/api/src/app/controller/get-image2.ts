import { assertRepositoryName, Repository } from '@lib/core';
import { Request, Response } from 'express';
import { injectable } from 'tsyringe';
import { ContributorsImageQuery } from '../query/contributors-image';
import { ContributorsRepository } from '../repository/contributors';
import { ContributorsImageRenderer } from '../service/contributors-image-renderer';
import { Controller } from '../utils/types';

@injectable()
export class GetImage2Controller implements Controller {
  constructor(
    private readonly imageQuery: ContributorsImageQuery,
    private readonly contributorsRepository: ContributorsRepository,
    private readonly renderer: ContributorsImageRenderer,
  ) {}
  async onRequest(req: Request, res: Response) {
    const repoName = req.query.repo;
    if (!assertRepositoryName(repoName)) {
      res.status(400).send(`"${repoName}" is not a valid repository name`);
      return;
    }
    const repository = Repository.fromString(repoName);
    const contributors = await this.contributorsRepository.getAllContributors(repository);
    const svg = await this.renderer.renderSvg(contributors);

    res.status(200).send(svg);
  }
}
