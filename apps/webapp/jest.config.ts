/* eslint-disable */
const path = require('path');
const { pathsToModuleNameMapper } = require('ts-jest');
const { compilerOptions } = require('../../tsconfig.base.json');

export default {
  preset: 'jest-preset-angular',
  globalSetup: 'jest-preset-angular/global-setup',
  displayName: 'webapp',
  moduleNameMapper: pathsToModuleNameMapper(compilerOptions.paths, { prefix: path.resolve(__dirname, '../../') }),
  coverageDirectory: '../../coverage/apps/webapp',
  setupFilesAfterEnv: ['<rootDir>/src/test-setup.ts'],
};
