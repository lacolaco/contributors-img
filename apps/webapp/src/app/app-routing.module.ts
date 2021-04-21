import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import { PreviewModule } from './features/preview/preview.module';

export const routes: Routes = [
  {
    path: '',
    pathMatch: 'full',
    redirectTo: 'preview',
  },
  {
    path: 'preview2',
    loadChildren: () => import('./features/preview2/preview2.module').then((m) => m.Preview2Module),
  },
];

@NgModule({
  imports: [RouterModule.forRoot(routes), PreviewModule],
  exports: [RouterModule],
})
export class AppRoutingModule {}
