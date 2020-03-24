import { createComponentFactory, Spectator } from '@ngneat/spectator';
import { SimpleRenderComponent } from './simple.component';
import { CONTRIBUTORS_DATA } from '@api/shared/state/tokens';
import { ContributorsListModule } from 'src/app/shared/contributors-list/contributors-list.module';

describe('SimpleRenderComponent', () => {
  let spectator: Spectator<SimpleRenderComponent>;
  const createComponent = createComponentFactory({
    component: SimpleRenderComponent,
    imports: [ContributorsListModule],
  });

  beforeEach(() => {
    (window as any)[CONTRIBUTORS_DATA] = [];
  });

  it('should create', () => {
    spectator = createComponent({});

    expect(spectator.component).toBeTruthy();
  });
});
