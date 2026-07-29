import { Injectable, signal } from '@angular/core';
import { ImageParams } from '../../models/image-params';
import { Repository } from '../../models/repository';

export interface State {
  imageParams: ImageParams;
  fetchingCount: number;
  result: {
    data: string;
  } | null;
}

export const defaultImageParams: ImageParams = {
  repository: new Repository('angular', 'angular-ja'),
  max: null,
  columns: null,
};

export const initialValue: State = {
  imageParams: defaultImageParams,
  fetchingCount: 0,
  result: null,
};

@Injectable()
export class PreviewState {
  readonly #state = signal<State>(initialValue);
  readonly state = this.#state.asReadonly();

  startFetchingImage() {
    this.#state.update((state) => ({
      ...state,
      fetchingCount: state.fetchingCount + 1,
      result: null,
    }));
  }

  finishFetchingImage(result: { data: string } | null) {
    this.#state.update((state) => ({
      ...state,
      fetchingCount: state.fetchingCount - 1,
      result,
    }));
  }

  patchImageParams(params: Partial<ImageParams>) {
    this.#state.update((state) => ({
      ...state,
      imageParams: {
        ...state.imageParams,
        ...params,
      },
    }));
  }
}
