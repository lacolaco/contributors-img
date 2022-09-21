import { Repository } from './repository';

export type ImageParamsFormValue = {
  repository: string | null;
  max: number | null;
  columns: number | null;
};

export type ImageParams = {
  repository: Repository;
  max: number | null;
  columns: number | null;
};
