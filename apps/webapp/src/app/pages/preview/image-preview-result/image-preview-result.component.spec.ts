import { ComponentFixture, TestBed } from '@angular/core/testing';

import { ImagePreviewResultComponent } from './image-preview-result.component';
import { Repository } from '../../../models';

describe('ImagePreviewResultComponent', () => {
  let component: ImagePreviewResultComponent;
  let fixture: ComponentFixture<ImagePreviewResultComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [ImagePreviewResultComponent],
    }).compileComponents();

    fixture = TestBed.createComponent(ImagePreviewResultComponent);
    fixture.componentRef.setInput('repository', Repository.fromString('foo/bar'));
    fixture.componentRef.setInput('imageSvg', 'foo');
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
