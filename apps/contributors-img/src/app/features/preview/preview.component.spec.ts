import { RouterTestingModule } from '@angular/router/testing';
import { createComponentFactory, Spectator } from '@ngneat/spectator';
import { PreviewComponent } from './preview.component';
import { PreviewModule } from './preview.module';

xdescribe('PreviewComponent', () => {
  let spectator: Spectator<PreviewComponent>;
  const createComponent = createComponentFactory({
    component: PreviewComponent,
    imports: [PreviewModule, RouterTestingModule],
  });

  it('should create', () => {
    spectator = createComponent({});

    expect(spectator.component).toBeTruthy();
  });
});
