import * as sharp from 'sharp';

export function convertToWebp(image: Buffer): Promise<Buffer> {
  return sharp(image).webp({}).toBuffer();
}
