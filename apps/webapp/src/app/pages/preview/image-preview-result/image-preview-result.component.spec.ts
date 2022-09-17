import { ComponentFixture, TestBed } from '@angular/core/testing';

import { ImagePreviewResultComponent } from './image-preview-result.component';

describe('ImagePreviewResultComponent', () => {
  let component: ImagePreviewResultComponent;
  let fixture: ComponentFixture<ImagePreviewResultComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [ImagePreviewResultComponent],
    }).compileComponents();

    fixture = TestBed.createComponent(ImagePreviewResultComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
