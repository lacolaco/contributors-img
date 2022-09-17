import { Routes } from '@angular/router';
import { PreviewPageComponent } from './pages/preview/preview.component';

export const routes: Routes = [
  {
    path: '',
    pathMatch: 'full',
    redirectTo: 'preview',
  },
  {
    path: 'preview',
    component: PreviewPageComponent,
  },
];
