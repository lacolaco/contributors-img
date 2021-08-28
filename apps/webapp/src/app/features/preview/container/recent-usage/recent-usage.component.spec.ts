import { ComponentFixture, TestBed } from '@angular/core/testing';
import { AngularFirestore } from '@angular/fire/compat/firestore';
import { of } from 'rxjs';
import { RepositoryGalleryComponent } from '../../component/repository-gallery/repository-gallery.component';
import { RecentUsageComponent } from './recent-usage.component';

describe('RecentUsageComponent', () => {
  let component: RecentUsageComponent;
  let fixture: ComponentFixture<RecentUsageComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [RecentUsageComponent, RepositoryGalleryComponent],
      providers: [
        {
          provide: AngularFirestore,
          useValue: {
            collection: () => ({
              valueChanges: jest.fn().mockReturnValue(of([])),
            }),
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
