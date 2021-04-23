import { Contributor } from '@lib/core';
import * as puppeteer from 'puppeteer';
import { injectable } from 'tsyringe';
import { environment } from '../../environments/environment';
import { createSvgInstance } from '../utils/svg-builder';

// eslint-disable-next-line @typescript-eslint/no-explicit-any
let rendererApp: any;

@injectable()
export class ContributorsImageRenderer {
  async render(contributors: Contributor[]): Promise<Buffer> {
    rendererApp = rendererApp ?? require('../../../../../dist/apps/renderer/server/main.js');

    const [html, browser] = await Promise.all([
      rendererApp.renderModule(rendererApp.AppServerModule, {
        document: '<app-root></app-root>',
        extraProviders: [rendererApp.provideContributors(contributors)],
      }),
      puppeteer.launch({
        headless: environment.useHeadless,
        args: ['--no-sandbox'],
      }),
    ]);

    const page = await browser.newPage();
    await page.setViewport({ width: 1048, height: 1048 });

    page.on('error', (error) => {
      console.error(error);
    });

    await page.setContent(html, { waitUntil: 'networkidle0' });
    const screenshotTarget = await page.waitForSelector('#renderTarget', { timeout: 0 });
    if (!screenshotTarget) {
      throw new Error('screenshot target element is not found.');
    }
    const screenshot = await screenshotTarget.screenshot({ type: 'png', omitBackground: true });
    if (!screenshot) {
      throw new Error('screenshot failed.');
    }

    return await browser.close().then(() => (typeof screenshot === 'string' ? Buffer.from(screenshot) : screenshot));
  }

  async renderSvg(contributors: Contributor[]) {
    const avatarSize = 64;
    const avatarGap = 4;
    const columnCount = 12;
    const rowCount = Math.ceil(contributors.length / columnCount);

    const container = createSvgInstance();
    container.size(
      (avatarSize + avatarGap) * (columnCount - 1) + avatarSize,
      (avatarSize + avatarGap) * (rowCount - 1) + avatarSize,
    );

    for (const [i, { avatar_url, html_url, login }] of contributors.entries()) {
      const x = (i % columnCount) * (avatarSize + avatarGap);
      const y = Math.floor(i / columnCount) * (avatarSize + avatarGap);

      const avatarCircle = container
        .circle(avatarSize)
        .stroke({
          color: '#c0c0c0',
          width: 1,
        })
        .fill(container.pattern(avatarSize, avatarSize).add(container.image(avatar_url).size(avatarSize, avatarSize)))
        .linkTo((link) => link.to(html_url).target('_blank'))
        .add(container.element('title').words(login));

      const inner = container.nested();
      inner.add(avatarCircle);
      inner.move(x, y);
    }

    return container.svg();
  }
}
