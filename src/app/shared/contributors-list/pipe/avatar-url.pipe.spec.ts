import { createPipeFactory, SpectatorPipe } from '@ngneat/spectator';

import { AvatarUrlPipe } from './avatar-url.pipe';

describe('AvatarUrlPipe ', () => {
  let spectator: SpectatorPipe<AvatarUrlPipe>;
  const createPipe = createPipeFactory(AvatarUrlPipe);

  it('should change the background color', () => {
    spectator = createPipe(`<div>{{ 'https://avatars0.githubusercontent.com/u/000000?v=4' | avatarUrl:64 }}</div>`);

    expect(spectator.element).toHaveText('https://avatars0.githubusercontent.com/u/000000?v=4&size=64');
  });
});
