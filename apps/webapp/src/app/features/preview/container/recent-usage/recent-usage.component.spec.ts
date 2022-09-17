import { ComponentFixture, TestBed } from '@angular/core/testing';
import { of } from 'rxjs';
import { FeaturedRepositoryDatastore } from '../../service/featured-repository-datastore';
import { RepositoryGalleryComponent } from '../../component/repository-gallery/repository-gallery.component';
import { RecentUsageComponent } from './recent-usage.component';

describe('RecentUsageComponent', () => {
  let component: RecentUsageComponent;
  let fixture: ComponentFixture<RecentUsageComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [RecentUsageComponent, RepositoryGalleryComponent],
      providers: [
        {
          provide: FeaturedRepositoryDatastore,
          useValue: {
            repositories$: of([]),
          },
        },
      ],
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
