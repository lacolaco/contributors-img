import 'core-js/features/reflect';
import 'zone.js/dist/zone-node';
import { environment } from './environments/environment';

if (environment.production) {
  // eslint-disable-next-line @typescript-eslint/no-var-requires
  require('@google-cloud/trace-agent').start();
}

import { startServer } from './app';

startServer({
  port: process.env.PORT || 3333,
});
