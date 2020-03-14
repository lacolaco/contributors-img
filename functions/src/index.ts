import * as functions from 'firebase-functions';
import * as admin from 'firebase-admin';
import * as puppeteer from 'puppeteer';
import * as Octokit from '@octokit/rest';
import * as cors from 'cors';

admin.initializeApp();

const octokit = new Octokit({
  auth: 'token 393ad1f410e7f6e6d78a19466812b6cea4d1ed52',
});
const withCors = cors({ origin: true });

const bucket = admin.storage().bucket();

function generateCacheId(repository: string) {
  return `image-cache--${repository.replace('/', '-')}`;
}

async function renderContributorsImage(repository: string): Promise<Buffer> {
  const browser = await puppeteer.launch({
    headless: true,
    args: ['--no-sandbox'],
  });
  const page = await browser.newPage();
  await page.setViewport({
    width: 1048,
    height: 1048,
  });

  console.log(`renderContributorsImage--Go to Page`);
  await page.goto(`https://contributors-img.firebaseapp.com?repo=${repository}`, { waitUntil: 'networkidle0' });

  console.log(`renderContributorsImage--Wait for selector`);
  const screenshotTarget = await page.waitForSelector('#contributors', { timeout: 0 });

  console.log(`renderContributorsImage--Wait for screenshot`);
  const screenshot = await screenshotTarget.screenshot({ type: 'png', omitBackground: true });
  return await browser.close().then(() => screenshot);
}

export const createContributorsImage = functions
  .runWith({
    timeoutSeconds: 60,
    memory: '1GB',
  })
  .https.onRequest(async (request, response) => {
    const repoParam = request.query['repo'];

    if (!repoParam || typeof repoParam !== 'string') {
      response.status(400).send(`'repo' parameter is required.`);
      return;
    }

    const cacheId = generateCacheId(repoParam);
    const useCache = request.hostname !== 'localhost';

    async function readFile(filename: string): Promise<Buffer | null> {
      console.log(`readFile: ${filename}`);
      const file = bucket.file(filename);
      return await file.exists().then(([exists]) => {
        if (exists) {
          return file.download({}).then(([data]) => data);
        }
        return null;
      });
    }

    async function writeFile(filename: string, file: Buffer): Promise<void> {
      console.log(`writeFile: ${filename}`);
      await bucket.file(filename).save(file, {
        public: true,
      });
    }

    const createImage = async () => {
      if (useCache) {
        console.debug('restore cache');
        const cache = await readFile(cacheId);
        if (cache) {
          console.debug('cache is found');
          return cache;
        }
      }
      console.log(`Render an image`);
      return renderContributorsImage(repoParam);
    };

    try {
      const image = await createImage();
      if (useCache) {
        console.debug('save cache');
        await writeFile(cacheId, image);
      }
      response
        .header('Content-Type', 'image/png')
        .header('Cache-Control', 'max-age=0, no-cache')
        .status(200)
        .send(image);
    } catch (error) {
      console.error(error);
      response.status(500).send(error);
    }
  });

export const getContributors = functions.https.onRequest((request, response) => {
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
});
