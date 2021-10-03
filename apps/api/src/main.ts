import 'core-js/features/reflect';
import { environment } from './environments/environment';
import { startServer } from './app';
import { container } from 'tsyringe';
import { provideInfrastructure } from './app/infrastructure/providers';

if (environment.production) {
  // eslint-disable-next-line @typescript-eslint/no-var-requires
  require('@google-cloud/trace-agent').start();
}

provideInfrastructure(container);

startServer({
  port: process.env.PORT || 3333,
});
