import { ChangeDetectionStrategy, Component, inject, Input } from '@angular/core';
import { FormBuilder, ReactiveFormsModule, Validators } from '@angular/forms';
import { MatButtonModule } from '@angular/material/button';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatInputModule } from '@angular/material/input';
import { ImageParams } from '../../../models/image-params';
import { Repository } from '../../../models/repository';
import { PreviewState } from '../state';

@Component({
  selector: 'app-image-preview-form',
  template: `
    <div class="controlPane">
      <form class="repositoryForm" [formGroup]="form">
        <mat-form-field appearance="outline" class="--no-hint" style="width: 320px">
          <mat-label>Enter GitHub Repository</mat-label>
          <input matInput name="repository" placeholder="e.g. angular/angular-ja" formControlName="repository" />
        </mat-form-field>
        <button type="submit" (click)="generateImage()" mat-raised-button color="primary" [disabled]="form.invalid">
          Generate
        </button>
      </form>
    </div>
  `,
  styleUrls: ['./image-preview-form.component.scss'],
  changeDetection: ChangeDetectionStrategy.OnPush,
  standalone: true,
  imports: [ReactiveFormsModule, MatFormFieldModule, MatInputModule, MatButtonModule],
})
export class ImagePreviewFormComponent {
  @Input()
  set value(value: ImageParams) {
    this.form.patchValue({ repository: value.repository.toString() });
  }

  private readonly store = inject(PreviewState);
  private readonly fb = inject(FormBuilder).nonNullable;

  readonly form = this.fb.group({
    repository: this.fb.control('', {
      validators: [Validators.required],
    }),
  });

  generateImage() {
    const repoName = this.form.getRawValue().repository;

    this.store.patchImageParams({
      repository: Repository.fromString(repoName),
    });
  }
}
