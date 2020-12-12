import 'core-js/features/reflect';

import { startServer } from './app';

startServer({
  port: process.env.PORT || 3333,
});
