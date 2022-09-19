import { RepositoryImageUrlPipe } from './repository-image-url.pipe';

describe('RepositoryOwnerPipe', () => {
  it('create an instance', () => {
    const pipe = new RepositoryImageUrlPipe();
    expect(pipe).toBeTruthy();
  });
});
