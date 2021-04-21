import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';

import { Preview2RoutingModule } from './preview2-routing.module';
import { Preview2Component } from './preview2.component';

@NgModule({
  declarations: [Preview2Component],
  imports: [CommonModule, Preview2RoutingModule],
})
export class Preview2Module {}
