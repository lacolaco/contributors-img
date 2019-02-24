import * as functions from 'firebase-functions';
import * as firebase from 'firebase-admin';
import * as puppeteer from 'puppeteer';

firebase.initializeApp();
const bucket = firebase.storage().bucket();
const isDebug = process.env.NODE_ENV !== 'production';

function generateCacheId(repository: string) {
  return `image-cache--${repository.replace('/', '-')}`;
}

async function renderContributorsImage(repository: string): Promise<Buffer> {
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
  return await browser.close().then(() => screenshot);
}

async function _createContributorsImage(repository: string): Promise<Buffer> {
  const cacheId = generateCacheId(repository);
  const cacheFile = bucket.file(cacheId);

  console.log(`Look for a cache...`);
  if (await cacheFile.exists().then(data => data[0])) {
    console.log(`Return from the cache`);
    return cacheFile.download().then(data => data[0]);
  }

  console.log(`Render an image`);
  const image = await renderContributorsImage(repository);

  console.log(`Save new cache`);
  await cacheFile.save(image, {});

  console.log(`Return rendered image`);
  return image;
}

export const createContributorsImage = functions
  .runWith({
    timeoutSeconds: 15,
    memory: '1GB',
  })
  .https.onRequest((request, response) => {
    const repo = request.param('repo');

    if (!repo || typeof repo !== 'string') {
      response.status(400).send(`'repo' parameter is required.`);
      return;
    }

    _createContributorsImage(repo)
      .then(image => {
        response.setHeader('Content-Type', 'image/png');
        response.status(200).send(image);
      })
      .catch(err => {
        console.error(err);
        response.status(500).send(err.toString());
      });
  });
