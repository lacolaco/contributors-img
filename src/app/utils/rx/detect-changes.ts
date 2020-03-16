import { ɵmarkDirty as markDirty } from '@angular/core';
import { pipe } from 'rxjs';
import { tap } from 'rxjs/operators';

export const scheduleChangeDetection = <T>(component: {}) =>
  pipe(
    tap<T>(() => markDirty(component)),
  );
