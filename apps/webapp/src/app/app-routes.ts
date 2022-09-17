import { Routes } from '@angular/router';
import { PreviewComponent } from './features/preview/preview.component';

export const routes: Routes = [
  {
    path: '',
    pathMatch: 'full',
    redirectTo: 'preview',
  },
  {
    path: 'preview',
    component: PreviewComponent,
  },
];
