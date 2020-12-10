import { async, ComponentFixture, TestBed } from '@angular/core/testing';
import { MatCardModule } from '@angular/material/card';
import { RouterTestingModule } from '@angular/router/testing';
import { ContributorsListModule } from '../../../../shared/contributors-list/contributors-list.module';
import { RepositoryGalleryComponent } from './repository-gallery.component';

describe('RepositoryGalleryComponent', () => {
  let component: RepositoryGalleryComponent;
  let fixture: ComponentFixture<RepositoryGalleryComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [RepositoryGalleryComponent],
      imports: [RouterTestingModule, MatCardModule, ContributorsListModule],
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
