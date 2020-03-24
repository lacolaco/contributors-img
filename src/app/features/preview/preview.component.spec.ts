import { createComponentFactory, Spectator } from '@ngneat/spectator';
import { PreviewComponent } from './preview.component';
import { PreviewModule } from './preview.module';

describe('PreviewComponent', () => {
  let spectator: Spectator<PreviewComponent>;
  const createComponent = createComponentFactory({
    component: PreviewComponent,
    imports: [PreviewModule],
  });

  it('should create', () => {
    spectator = createComponent({});

    expect(spectator.component).toBeTruthy();
  });
});
