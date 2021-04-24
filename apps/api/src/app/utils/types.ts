import { Request, Response } from 'express';

export interface Controller {
  onRequest(req: Request, res: Response): void;
}

export type SupportedImageType = 'image/png' | 'image/webp' | 'image/svg+xml';
