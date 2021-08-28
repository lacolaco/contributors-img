import { ComponentFixture, TestBed } from '@angular/core/testing';
import { AngularFireModule } from '@angular/fire/compat';
import { AngularFirestoreModule } from '@angular/fire/compat/firestore';
import { RecentUsageComponent } from './recent-usage.component';

describe('RecentUsageComponent', () => {
  let component: RecentUsageComponent;
  let fixture: ComponentFixture<RecentUsageComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [RecentUsageComponent],
      imports: [
        AngularFireModule.initializeApp({
          projectId: 'test',
        }),
        AngularFirestoreModule,
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
