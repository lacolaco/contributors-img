import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { SimpleRenderComponent } from './simple.component';

describe('SimpleComponent', () => {
  let component: SimpleRenderComponent;
  let fixture: ComponentFixture<SimpleRenderComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [SimpleRenderComponent],
    }).compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(SimpleRenderComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
