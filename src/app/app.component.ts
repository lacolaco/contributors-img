import { Component, ChangeDetectorRef } from '@angular/core';

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.scss'],
})
export class AppComponent {
  constructor(private cdRef: ChangeDetectorRef) {}

  onRouteActivate() {
    // this.cdRef.detectChanges();
  }
}
