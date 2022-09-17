import { ComponentFixture, TestBed } from '@angular/core/testing';
import { provideNoopFeaturedRepositoryDatasource } from '../../../shared/featured-repository/noop';
import { RepositoryGalleryComponent } from '../repository-gallery/repository-gallery.component';
import { RecentUsageComponent } from './recent-usage.component';

describe('RecentUsageComponent', () => {
  let component: RecentUsageComponent;
  let fixture: ComponentFixture<RecentUsageComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [RecentUsageComponent, RepositoryGalleryComponent],
      providers: [provideNoopFeaturedRepositoryDatasource()],
    }).compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(RecentUsageComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
