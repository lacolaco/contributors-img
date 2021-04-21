import { ComponentFixture, TestBed } from '@angular/core/testing';

import { Preview2Component } from './preview2.component';

describe('Preview2Component', () => {
  let component: Preview2Component;
  let fixture: ComponentFixture<Preview2Component>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [Preview2Component],
    }).compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(Preview2Component);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
