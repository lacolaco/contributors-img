import { NO_ERRORS_SCHEMA } from '@angular/core';
import { ReactiveFormsModule } from '@angular/forms';
import { createComponentFactory, Spectator } from '@ngneat/spectator';
import { Repository } from '@api/shared/model/repository';
import { FormComponent } from './form.component';

describe('FormComponent', () => {
  let spectator: Spectator<FormComponent>;
  const createComponent = createComponentFactory({
    component: FormComponent,
    imports: [ReactiveFormsModule],
    schemas: [NO_ERRORS_SCHEMA],
  });

  beforeEach(() => {});

  it('should create', () => {
    spectator = createComponent({});

    expect(spectator.component).toBeTruthy();
  });

  it('should set form value on input update', () => {
    spectator = createComponent({});
    spectator.setInput({
      repository: new Repository('foo', 'bar'),
    });

    expect(spectator.component.form.value).toEqual({
      repository: 'foo/bar',
    });
  });
});
