import * as puppeteer from 'puppeteer';
import { environment } from '../../environment';

export async function useBrowserPage<T>(
  options: { viewport: { width: number; height: number } },
  callback: (page: puppeteer.Page) => Promise<T>,
): Promise<T> {
  const browser = await puppeteer.launch({
    args: ['--no-sandbox'],
    headless: environment.production,
  });

  const page = await browser.newPage();
  await page.setViewport(options.viewport);

  const result = await callback(page);
  await browser.close();
  return result;
}
