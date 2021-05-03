import { ComponentFixture, TestBed } from '@angular/core/testing';

import { SvgViewComponent } from './svg-view.component';

describe('SvgViewComponent', () => {
  let component: SvgViewComponent;
  let fixture: ComponentFixture<SvgViewComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [SvgViewComponent],
    }).compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(SvgViewComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
