import { NO_ERRORS_SCHEMA } from '@angular/core';
import { ReactiveFormsModule } from '@angular/forms';
import { Repository } from '@lib/core';
import { createComponentFactory, Spectator } from '@ngneat/spectator';
import { RepositoryFormComponent } from './repository-form.component';

describe('RepositoryFormComponent', () => {
  let spectator: Spectator<RepositoryFormComponent>;
  const createComponent = createComponentFactory({
    component: RepositoryFormComponent,
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
