import { ComponentFixture, TestBed } from '@angular/core/testing';
import { PreviewState } from '../state';

import { ImagePreviewComponent } from './image-preview.component';
import { provideNoopAnimations } from '@angular/platform-browser/animations';

describe('ImagePreviewComponent', () => {
  let component: ImagePreviewComponent;
  let fixture: ComponentFixture<ImagePreviewComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [ImagePreviewComponent],
      providers: [PreviewState, provideNoopAnimations()],
    }).compileComponents();

    fixture = TestBed.createComponent(ImagePreviewComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
