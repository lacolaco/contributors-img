import { Repository } from '../model/repository';
import { useBrowserPage } from './utils/use-browser-page';

export async function generateContributorsImage(repository: Repository): Promise<Buffer> {
  return await useBrowserPage(
    {
      viewport: {
        width: 1048,
        height: 1048,
      },
    },
    async page => {
      console.debug(`open app`);
      await page.goto(`https://contributors-img.firebaseapp.com?repo=${repository.toString()}`, {
        waitUntil: 'networkidle0',
      });

      console.debug(`wait for schreenshot target`);
      const targetElement = await page.waitForSelector('#contributors', {
        timeout: 0, // disable timeout
      });

      console.debug(`take screenshot`);
      return await targetElement.screenshot({ type: 'png', omitBackground: true });
    },
  );
}
