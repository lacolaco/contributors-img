import { NO_ERRORS_SCHEMA } from '@angular/core';
import { ReactiveFormsModule } from '@angular/forms';
import { createComponentFactory, Spectator } from '@ngneat/spectator';
import { Repository } from 'shared/model/repository';
import { AppStore } from '../state/store';
import { FormComponent } from './form.component';

describe('FormComponent', () => {
  let spectator: Spectator<FormComponent>;
  let store: AppStore;
  const createComponent = createComponentFactory({
    component: FormComponent,
    imports: [ReactiveFormsModule],
    schemas: [NO_ERRORS_SCHEMA],
  });

  beforeEach(() => {});

  it('should create', () => {
    spectator = createComponent({});
    store = spectator.inject(AppStore);

    expect(spectator.component).toBeTruthy();
  });

  it('should set form value on store update', () => {
    spectator = createComponent({});
    store = spectator.inject(AppStore);
    store.update(state => ({ ...state, repository: new Repository('foo', 'bar') }));
    spectator.detectChanges();

    expect(spectator.component.form.value).toEqual({
      repository: 'foo/bar',
    });
  });
});
