import 'jest-preset-angular/setup-jest';

import { getTestBed } from '@angular/core/testing';
import { BrowserDynamicTestingModule, platformBrowserDynamicTesting } from '@angular/platform-browser-dynamic/testing';
import { NoopAnimationsModule } from '@angular/platform-browser/animations';

getTestBed().resetTestEnvironment();
getTestBed().initTestEnvironment([BrowserDynamicTestingModule, NoopAnimationsModule], platformBrowserDynamicTesting(), {
  teardown: { destroyAfterEach: true },
});
