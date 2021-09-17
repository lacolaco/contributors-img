import { Request, Response } from 'express';
import { Readable } from 'stream';

export interface Controller {
  onRequest(req: Request, res: Response): void;
}

export interface FileStream {
  readonly data: Readable;
  readonly contentType: string;
}
