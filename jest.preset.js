const nxPreset = require('@nrwl/jest/preset');

module.exports = {
  ...nxPreset,
  transform: {
    '^.+\\.(ts|js|html)$': 'ts-jest',
  },
  resolver: '@nrwl/jest/plugins/resolver',
  moduleFileExtensions: ['ts', 'js', 'html'],
};
