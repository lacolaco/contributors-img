/* eslint-disable */
const presetTs = require('ts-jest/presets').defaults;

export default {
  ...presetTs,
  displayName: 'core',
  globals: {
    'ts-jest': {
      tsconfig: '<rootDir>/tsconfig.spec.json',
    },
  },
  coverageDirectory: '../../coverage/libs/core',
};
