import fetch from 'node-fetch';
import AbortController from 'abort-controller';

export async function createDataURIFromURL(imageUrl: string): Promise<string> {
  try {
    const controller = new AbortController();
    setTimeout(() => {
      controller.abort();
    }, 30 * 1000);
    const resp = await fetch(imageUrl, {
      signal: controller.signal,
    });
    const contentType = resp.headers.get('content-type');
    const data = (await resp.buffer()).toString('base64');
    return `data:${contentType};base64,${data}`;
  } catch (error) {
    console.error(error);
    return '';
  }
}
