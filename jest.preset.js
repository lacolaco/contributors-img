const { pathsToModuleNameMapper } = require('ts-jest');
const { compilerOptions } = require('./tsconfig.base.json');

module.exports = {
  globalSetup: 'jest-preset-angular/global-setup',
  moduleNameMapper: pathsToModuleNameMapper(compilerOptions.paths, { prefix: __dirname }),
};
