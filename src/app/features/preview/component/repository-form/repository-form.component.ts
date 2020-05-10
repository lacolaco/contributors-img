import { Component, EventEmitter, Input, Output, ChangeDetectionStrategy } from '@angular/core';
import { FormControl, FormGroup, Validators } from '@angular/forms';
import { Repository } from '@api/shared/model';

@Component({
  selector: 'app-repository-form',
  templateUrl: './repository-form.component.html',
  styleUrls: ['./repository-form.component.scss'],
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class RepositoryFormComponent {
  @Input()
  set repository(value: Repository | null) {
    this.form.patchValue({ repository: value ? value.toString() : null });
  }

  @Output() valueChange = new EventEmitter<string>();

  readonly form = new FormGroup({
    repository: new FormControl(this.repository, {
      validators: [Validators.required],
    }),
  });

  generateImage() {
    const repoName = this.form.value.repository;
    this.valueChange.emit(repoName);
  }
}
