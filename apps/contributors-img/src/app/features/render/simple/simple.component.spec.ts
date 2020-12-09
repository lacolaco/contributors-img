import { CONTRIBUTORS_DATA } from '@api/shared/state/tokens';
import { createComponentFactory, Spectator } from '@ngneat/spectator';
import { ContributorsListModule } from '../../../shared/contributors-list/contributors-list.module';
import { SimpleRenderComponent } from './simple.component';

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
