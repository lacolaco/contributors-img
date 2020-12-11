import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';

import { RenderRoutingModule } from './render-routing.module';
import { SimpleRenderComponent } from './simple/simple.component';
import { ContributorsListModule } from '../../shared/contributors-list/contributors-list.module';

@NgModule({
  declarations: [SimpleRenderComponent],
  imports: [CommonModule, RenderRoutingModule, ContributorsListModule],
})
export class RenderModule {}
