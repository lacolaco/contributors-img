import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import { SimpleRenderComponent } from './features/render/simple/simple.component';
import { RenderModule } from './features/render/render.module';

export const routes: Routes = [
  {
    path: '',
    pathMatch: 'full',
    redirectTo: 'preview',
  },
  {
    path: 'preview',
    loadChildren: () => import('./features/preview/preview.module').then((m) => m.PreviewModule),
  },
  {
    path: 'render',
    children: [
      {
        path: 'simple',
        component: SimpleRenderComponent,
      },
    ],
  },
];

@NgModule({
  imports: [RouterModule.forRoot(routes, { relativeLinkResolution: 'legacy' }), RenderModule],
  exports: [RouterModule],
})
export class AppRoutingModule {}
