import { Logging } from '@google-cloud/logging';

export const createLogging = (): Logging => {
  return new Logging();
};
