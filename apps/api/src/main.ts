import * as tracing from '@google-cloud/trace-agent';
import 'core-js/features/reflect';
import { startServer } from './app';
import { environment } from './environments/environment';

if (environment.production) {
  tracing.start();
}

startServer({
  port: process.env.PORT || 3333,
});
