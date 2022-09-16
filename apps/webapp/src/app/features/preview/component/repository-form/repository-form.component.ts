import { ChangeDetectionStrategy, Component, EventEmitter, Input, Output } from '@angular/core';
import { FormControl, FormGroup, Validators } from '@angular/forms';
import { Repository } from '../../../../shared/model/repository';

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
    repository: new FormControl('', {
      validators: [Validators.required],
    }),
  });

  generateImage() {
    const repoName = this.form.value.repository;
    this.valueChange.emit(repoName ?? '');
  }
}
