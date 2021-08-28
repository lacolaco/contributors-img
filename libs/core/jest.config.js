const presetTs = require('ts-jest/presets').defaults;

module.exports = {
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
