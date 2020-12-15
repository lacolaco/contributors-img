import { InjectionToken } from '@angular/core';
import { Contributor } from '@lib/core';

export const CONTRIBUTORS_DATA = new InjectionToken<Contributor[]>('@@CONTRIBUTORS_DATA@@');
