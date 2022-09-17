import { ComponentFixture, TestBed } from '@angular/core/testing';
import { PreviewStore } from '../preview.store';

import { ImagePreviewComponent } from './image-preview.component';

describe('ImagePreviewComponent', () => {
  let component: ImagePreviewComponent;
  let fixture: ComponentFixture<ImagePreviewComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [ImagePreviewComponent],
      providers: [PreviewStore],
    }).compileComponents();

    fixture = TestBed.createComponent(ImagePreviewComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
