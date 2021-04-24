import fetch from 'node-fetch';

export async function createDataURIFromURL(imageUrl: string): Promise<string> {
  try {
    const resp = await fetch(imageUrl);
    const contentType = resp.headers.get('content-type');
    const data = (await resp.buffer()).toString('base64');
    return `data:${contentType};base64,${data}`;
  } catch (error) {
    console.error(error);
    return '';
  }
}
