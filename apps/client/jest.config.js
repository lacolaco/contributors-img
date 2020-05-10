module.exports = {
  name: 'contributors-img',
  preset: '../../jest.config.js',
  coverageDirectory: '../../coverage/apps/contributors-img',
  snapshotSerializers: [
    'jest-preset-angular/build/AngularNoNgAttributesSnapshotSerializer.js',
    'jest-preset-angular/build/AngularSnapshotSerializer.js',
    'jest-preset-angular/build/HTMLCommentSerializer.js'
  ]
};
