import { Component, EventEmitter, Input, Output, ChangeDetectionStrategy } from '@angular/core';
import { FormControl, FormGroup, Validators } from '@angular/forms';
import { Repository } from 'shared/model/repository';

@Component({
  selector: 'app-form',
  templateUrl: './form.component.html',
  styleUrls: ['./form.component.scss'],
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class FormComponent {
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
