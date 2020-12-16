import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import { PreviewModule } from './features/preview/preview.module';

export const routes: Routes = [
  {
    path: '',
    pathMatch: 'full',
    redirectTo: 'preview',
  },
];

@NgModule({
  imports: [RouterModule.forRoot(routes), PreviewModule],
  exports: [RouterModule],
})
export class AppRoutingModule {}
