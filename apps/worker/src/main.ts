import 'core-js/features/reflect';
import { environment } from './environments/environment';
import { startServer } from './app';

if (environment.production) {
  // eslint-disable-next-line @typescript-eslint/no-var-requires
  require('@google-cloud/trace-agent').start();
}

startServer({
  port: process.env.PORT || 3333,
});
