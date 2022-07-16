/* eslint-disable */
const presetAngular = require('jest-preset-angular/jest-preset');

export default {
  ...presetAngular,
  displayName: 'webapp',
  preset: '../../jest.preset.js',
  coverageDirectory: '../../coverage/apps/webapp',
  setupFilesAfterEnv: ['<rootDir>/src/test-setup.ts'],
};
