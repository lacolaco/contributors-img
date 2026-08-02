import { Component, ChangeDetectionStrategy } from '@angular/core';
import { RouterOutlet } from '@angular/router';

@Component({
  selector: 'app-root',
  template: `<router-outlet (activate)="onRouteActivate()" />`,
  styleUrls: ['./app.component.scss'],
  // Angular 22 made OnPush the default and its migration wrote this out to preserve the previous behaviour.
  // Switching to OnPush is a change in behaviour, not part of the upgrade, so it is deferred.
  // eslint-disable-next-line @angular-eslint/prefer-on-push-component-change-detection
  changeDetection: ChangeDetectionStrategy.Eager,
  imports: [RouterOutlet],
})
export class AppComponent {
  onRouteActivate() {
    // this.cdRef.detectChanges();
  }
}
