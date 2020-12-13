import * as compression from 'compression';
import * as express from 'express';
import * as morgan from 'morgan';
import routes from './routes';
import './infrastructure/providers';

export function startServer(config: { port: string | number }) {
  const app = express();
  app.use(compression({}));
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  app.use(morgan('combined') as any); // https://github.com/DefinitelyTyped/DefinitelyTyped/issues/50076
  app.use(routes());

  const server = app.listen(config.port, () => {
    console.log(`Listening at http://localhost:${config.port}`);
  });

  server.on('error', console.error);
  return server;
}
