/* eslint-disable @typescript-eslint/no-var-requires */
import { Request, Response } from 'express';
import { injectable } from 'tsyringe';
import { ContributorsImageRenderer } from '../service/contributors-image-renderer';
import { Controller } from '../utils/types';

@injectable()
export class RenderImageController implements Controller {
  constructor(private imageRenderer: ContributorsImageRenderer) {}
  async onRequest(req: Request, res: Response) {
    // eslint-disable-next-line @typescript-eslint/ban-ts-comment
    // @ts-ignore
    const { renderModule, AppServerModule } = await import('../../../../../dist/apps/renderer/server/main.js');

    const html = await renderModule(AppServerModule, {
      document: '<app-root></app-root>',
      // url: '/',
    });

    const image = await this.imageRenderer.render2(html);

    res.header('Content-Type', 'image/png').send(image);
  }
}
