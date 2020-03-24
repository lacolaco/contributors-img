import { Component, OnInit, Input } from '@angular/core';

@Component({
  selector: 'app-repository-gallery',
  templateUrl: './repository-gallery.component.html',
  styleUrls: ['./repository-gallery.component.scss'],
})
export class RepositoryGalleryComponent implements OnInit {
  constructor() {}

  @Input()
  repositories: any[];

  ngOnInit(): void {}
}
