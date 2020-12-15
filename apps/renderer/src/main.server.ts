import 'zone.js/dist/zone-node';

import { enableProdMode, Provider } from '@angular/core';
import { Contributor } from '@lib/core';
import { environment } from './environments/environment';
import { CONTRIBUTORS_DATA } from './app/tokens';

if (environment.production) {
  enableProdMode();
}

export function provideContributors(contributors: Contributor[]): Provider[] {
  return [{ provide: CONTRIBUTORS_DATA, useValue: contributors }];
}

export { renderModule, renderModuleFactory } from '@angular/platform-server';
export { AppServerModule } from './app/app.server.module';
