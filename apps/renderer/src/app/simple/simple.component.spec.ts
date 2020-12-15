import { ContributorsListModule } from '@lib/renderer-ui';
import { createComponentFactory, Spectator } from '@ngneat/spectator/jest';
import { SimpleRendererComponent } from './simple.component';

describe('SimpleRenderComponent', () => {
  let spectator: Spectator<SimpleRendererComponent>;
  const createComponent = createComponentFactory({
    component: SimpleRendererComponent,
    imports: [ContributorsListModule],
  });

  beforeEach(() => {
  });

  it('should create', () => {
    spectator = createComponent({});

    expect(spectator.component).toBeTruthy();
  });
});
