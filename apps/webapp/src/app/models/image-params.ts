import { Repository } from './repository';

export interface ImageParamsFormValue {
  repository: string | null;
  max: number | null;
  columns: number | null;
}

export interface ImageParams {
  repository: Repository;
  max: number | null;
  columns: number | null;
}
