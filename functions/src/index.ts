import * as functions from 'firebase-functions';
import * as firebase from 'firebase-admin';
import * as puppeteer from 'puppeteer';
import * as Octokit from '@octokit/rest';
import * as cors from 'cors';

firebase.initializeApp();
const octokit = new Octokit({
  auth: '393ad1f410e7f6e6d78a19466812b6cea4d1ed52',
});
const withCors = cors({ origin: true });

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

  await page.goto(
    `https://contributors-img.firebaseapp.com?repo=${repository}`,
  );

  const screenshotTarget = await page.waitForSelector('#contributors');

  await page.waitFor(10000);

  const screenshot = await screenshotTarget.screenshot({
    type: 'png',
    omitBackground: true,
  });
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
    timeoutSeconds: 30,
    memory: '1GB',
  })
  .https.onRequest((request, response) => {
    const repoParam = request.query['repo'];

    if (!repoParam || typeof repoParam !== 'string') {
      response.status(400).send(`'repo' parameter is required.`);
      return;
    }

    _createContributorsImage(repoParam)
      .then(image => {
        response.setHeader('Content-Type', 'image/png');
        response.status(200).send(image);
      })
      .catch(err => {
        console.error(err);
        response.status(500).send(err.toString());
      });
  });

export const getContributors = functions.https.onRequest(
  (request, response) => {
    withCors(request, response, () => {
      const repoParam = request.query['repo'];

      if (!repoParam || typeof repoParam !== 'string') {
        response.status(400).send(`'repo' parameter is required.`);
        return;
      }

      const [owner, repo] = repoParam.split('/');

      const options = octokit.repos.listContributors.endpoint.merge({
        owner,
        repo,
        per_page: 100,
      });
      octokit
        .paginate(options)
        .then(contributors => {
          response.setHeader('Content-Type', 'application/json');
          response.status(200).send(contributors);
        })
        .catch(err => {
          console.error(err);
          response.status(500).send(err.toString());
        });
    });
  },
);
