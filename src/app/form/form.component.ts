import { Component, OnInit, OnDestroy } from '@angular/core';
import { FormControl, FormGroup, Validators } from '@angular/forms';
import { Subject } from 'rxjs';
import { takeUntil } from 'rxjs/operators';
import { ContributorsStore } from '../state/contributors';

@Component({
  selector: 'app-form',
  templateUrl: './form.component.html',
  styleUrls: ['./form.component.scss'],
})
export class FormComponent implements OnInit, OnDestroy {
  form: FormGroup;

  private onDestroy$ = new Subject();

  constructor(private contributorsStore: ContributorsStore) {
    this.form = new FormGroup({
      repository: new FormControl(this.contributorsStore.value.repository, {
        validators: [Validators.required],
      }),
    });
  }

  ngOnInit() {
    this.contributorsStore
      .selectValue(state => state.repository)
      .pipe(takeUntil(this.onDestroy$))
      .subscribe(repository => {
        this.form.patchValue({ repository });
      });
  }

  generateImage() {
    const repository = this.form.value.repository;
    this.contributorsStore.updateValue(state => ({
      ...state,
      repository,
    }));
  }

  ngOnDestroy() {
    this.onDestroy$.next();
  }
}
