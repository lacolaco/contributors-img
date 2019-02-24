import * as functions from 'firebase-functions';
import * as firebase from 'firebase-admin';
import * as puppeteer from 'puppeteer';

firebase.initializeApp();
const bucket = firebase.storage().bucket();
const isDebug = process.env.NODE_ENV !== 'production';

function generateCacheId(repository: string) {
  return `image-cache--${repository.replace('/', '-')}`;
}

async function createScreenshot(repository: string): Promise<Buffer> {
  const cacheId = generateCacheId(repository);

  const cacheFile = bucket.file(cacheId);

  console.log(11);
  if (await cacheFile.exists().then(data => data[0])) {
    console.log(12);
    return cacheFile.download().then(data => data[0]);
  }
  console.log(22);

  const browser = await puppeteer.launch(
    isDebug
      ? {}
      : {
          headless: true,
          args: ['--no-sandbox'],
        },
  );

  const page = await browser.newPage();

  await page.goto(`https://github.com/${repository}/graphs/contributors`);

  await page.waitForResponse(
    `https://github.com/${repository}/graphs/contributors-data`,
  );
  const screenshotTarget = await page.waitForSelector('#contributors');

  const screenshot = await screenshotTarget.screenshot({ type: 'png' });

  console.log(33);

  await cacheFile.save(screenshot, {});

  console.log(44);

  return await browser.close().then(() => screenshot);
}

export const createContributorsImage = functions
  .runWith({
    timeoutSeconds: 15,
    memory: '1GB',
  })
  .https.onRequest((request, response) => {
    createScreenshot('angular/angular-ja')
      .then(image => {
        response.setHeader('Content-Type', 'image/png');
        response.status(200).send(image);
      })
      .catch(err => {
        console.error(err);
        response.status(500).send(err);
      });
  });
