import * as puppeteer from 'puppeteer';
import { Contributor, CONTRIBUTORS_DATA } from '@contributors-img/api-interfaces';

export async function renderContributorsImage({
  contributors,
  config,
}: {
  contributors: Contributor[];
  config: {
    webappUrl: string;
    useHeadless: boolean;
  };
}): Promise<Buffer> {
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

  page.on('error', (error) => {
    console.error(error);
  });

  await page.evaluateOnNewDocument(
    (token, json) => {
      (window as any)[token] = JSON.parse(json);
    },
    CONTRIBUTORS_DATA,
    JSON.stringify(contributors),
  );

  console.log(`go to webapp`);
  await page.goto(`${config.webappUrl}/render/simple`, { waitUntil: 'networkidle0' });

  console.log(`wait for selector`);
  const screenshotTarget = await page.waitForSelector('#renderTarget', { timeout: 0 });

  console.log(`capture screenshot`);
  const screenshot = await screenshotTarget.screenshot({ type: 'png', omitBackground: true });
  console.groupEnd();

  return await browser.close().then(() => screenshot);
}
