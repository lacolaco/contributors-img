import { Request, Response } from 'express';

export interface Controller {
  onRequest(req: Request, res: Response): void;
}
