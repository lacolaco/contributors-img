import { NO_ERRORS_SCHEMA } from '@angular/core';
import { RouterTestingModule } from '@angular/router/testing';
import { createComponentFactory, Spectator } from '@ngneat/spectator/jest';
import { ReactiveComponentModule } from '@ngrx/component';
import { PreviewComponent } from './preview.component';

xdescribe('PreviewComponent', () => {
  let spectator: Spectator<PreviewComponent>;
  const createComponent = createComponentFactory({
    component: PreviewComponent,
    schemas: [NO_ERRORS_SCHEMA],
    imports: [RouterTestingModule, ReactiveComponentModule],
  });

  test('should create', () => {
    spectator = createComponent({});

    expect(spectator.component).toBeTruthy();
  });
});
