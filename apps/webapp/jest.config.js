const presetAngular = require('jest-preset-angular/jest-preset');

module.exports = {
  ...presetAngular,
  displayName: 'webapp',
  preset: '../../jest.preset.js',
  coverageDirectory: '../../coverage/apps/webapp',
  setupFilesAfterEnv: ['<rootDir>/src/test-setup.ts'],
};
