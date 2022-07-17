/* eslint-disable */
const presetTs = require('ts-jest/presets').defaults;

export default {
  ...presetTs,
  displayName: 'core',
  preset: '../../jest.preset.js',
  globals: {
    'ts-jest': {
      tsconfig: '<rootDir>/tsconfig.spec.json',
    },
  },
  coverageDirectory: '../../coverage/libs/core',
};
