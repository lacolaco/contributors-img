import { ɵdetectChanges as detectChanges } from '@angular/core';
import { pipe } from 'rxjs';
import { tap } from 'rxjs/operators';

export const requestChangeDetection = <T>(component: {}) =>
  pipe(
    tap<T>(() => {
      requestAnimationFrame(() => detectChanges(component));
    }),
  );
