import { Request } from 'express';

export function isGitHubRequest(req: Request): boolean {
  const userAgent = req.headers['user-agent'] ?? '';
  return userAgent.toLowerCase().includes('github');
}
