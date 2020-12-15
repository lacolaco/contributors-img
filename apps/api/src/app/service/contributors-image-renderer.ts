import { Contributor, CONTRIBUTORS_DATA } from '@lib/core';
import * as puppeteer from 'puppeteer';
import { injectable } from 'tsyringe';
import { environment } from '../../environments/environment';

@injectable()
export class ContributorsImageRenderer {
  async render(contributors: Contributor[]): Promise<Buffer> {
    const browser = await puppeteer.launch({
      headless: environment.useHeadless,
      args: ['--no-sandbox'],
    });

    const page = await browser.newPage();
    await page.setViewport({ width: 1048, height: 1048 });

    page.on('error', (error) => {
      console.error(error);
    });

    await page.evaluateOnNewDocument(
      (token, json) => {
        // eslint-disable-next-line @typescript-eslint/no-explicit-any
        (window as any)[token] = JSON.parse(json);
      },
      CONTRIBUTORS_DATA,
      JSON.stringify(contributors),
    );

    await page.goto(`${environment.webappUrl}/render/simple`, { waitUntil: 'networkidle0' });
    const screenshotTarget = await page.waitForSelector('#renderTarget', { timeout: 0 });
    const screenshot = await screenshotTarget.screenshot({ type: 'png', omitBackground: true });

    return await browser.close().then(() => screenshot);
  }

  async render2(html: string): Promise<Buffer> {
    const browser = await puppeteer.launch({
      headless: environment.useHeadless,
      args: ['--no-sandbox'],
    });

    const page = await browser.newPage();
    await page.setViewport({ width: 1048, height: 1048 });

    page.on('error', (error) => {
      console.error(error);
    });

    await page.setContent(html, { waitUntil: 'networkidle0' });
    // const screenshotTarget = await page.waitForSelector('#renderTarget', { timeout: 0 });
    const screenshot = await page.screenshot({ type: 'png', omitBackground: true });

    return await browser.close().then(() => screenshot);
  }
}
