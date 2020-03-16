import * as puppeteer from 'puppeteer';
import { Repository } from '../shared/model/repository';

export async function renderContributorsImage(
  repository: Repository,
  config: { webappUrl: string; useHeadless: boolean },
): Promise<Buffer> {
  console.group('renderContributorsImage');

  console.log(`open browser`);
  const browser = await puppeteer.launch({
    headless: config.useHeadless,
    args: ['--no-sandbox'],
  });

  console.log(`open page`);
  const page = await browser.newPage();
  await page.setViewport({
    width: 1048,
    height: 1048,
  });

  console.log(`go to webapp`);
  await page.goto(`${config.webappUrl}?repo=${repository.toString()}`, { waitUntil: 'networkidle0' });

  console.log(`wait for selector`);
  const screenshotTarget = await page.waitForSelector('#contributors', { timeout: 0 });

  console.log(`capture screenshot`);
  const screenshot = await screenshotTarget.screenshot({ type: 'png', omitBackground: true });
  console.groupEnd();
  return await browser.close().then(() => screenshot);
}
