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
    const avatarCircleSize = 64;
    const avatarCircleGap = 4;
    const columnCount = 12;
    const rowCount = Math.ceil(contributors.length / columnCount);

    const svg = createSvgInstance();
    svg.size(
      (avatarCircleSize + avatarCircleGap) * (columnCount - 1) + avatarCircleSize,
      (avatarCircleSize + avatarCircleGap) * (rowCount - 1) + avatarCircleSize,
    );

    for (const [i, { avatar_url, html_url, login }] of contributors.entries()) {
      const x = (i % columnCount) * (avatarCircleSize + avatarCircleGap);
      const y = Math.floor(i / columnCount) * (avatarCircleSize + avatarCircleGap);
      const image = svg.image(avatar_url).size(avatarCircleSize, avatarCircleSize);
      const pattern = svg.pattern(avatarCircleSize, avatarCircleSize).move(x, y).add(image);
      svg
        .circle(avatarCircleSize)
        .fill(pattern)
        .stroke('#c0c0c0')
        .add(svg.element('title').words(login))
        .linkTo((link) => link.to(html_url).target('_blank'))
        .move(x, y);
    }

    return svg.svg();
  }
}
