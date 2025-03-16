import { Component } from '@angular/core';
import { RouterOutlet } from '@angular/router';

@Component({
  selector: 'app-root',
  template: `<router-outlet (activate)="onRouteActivate()"></router-outlet>`,
  styleUrls: ['./app.component.scss'],
  imports: [RouterOutlet],
})
export class AppComponent {
  onRouteActivate() {
    // this.cdRef.detectChanges();
  }
}
