import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { RepositoryGalleryComponent } from './repository-gallery.component';

describe('RepositoryGalleryComponent', () => {
  let component: RepositoryGalleryComponent;
  let fixture: ComponentFixture<RepositoryGalleryComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [RepositoryGalleryComponent],
    }).compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(RepositoryGalleryComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
