import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import { Preview2Component } from './preview2.component';

const routes: Routes = [
  {
    path: '',
    pathMatch: 'full',
    component: Preview2Component,
  },
];

@NgModule({
  imports: [RouterModule.forChild(routes)],
  exports: [RouterModule],
})
export class Preview2RoutingModule {}
