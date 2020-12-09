import { Component, OnInit } from '@angular/core';
import { Contributor } from '@api/shared/model';
import { CONTRIBUTORS_DATA } from '@api/shared/state/tokens';

@Component({
  selector: 'app-simple',
  templateUrl: './simple.component.html',
  styleUrls: ['./simple.component.scss'],
})
export class SimpleRenderComponent implements OnInit {
  contributors: Contributor[];

  constructor() {}

  ngOnInit(): void {
    const contributors = (window as any)[CONTRIBUTORS_DATA] as Contributor[];
    if (!contributors) {
      throw new Error('No Embedded Contributors Data');
    }
    this.contributors = contributors;
  }
}
