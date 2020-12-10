import { environment } from '../../environments/environment';

export interface ApplicationConfig {
  webappUrl: string;
  useHeadless: boolean;
  useCache: boolean;
}

export function getApplicationConfig(): ApplicationConfig {
  return {
    webappUrl: environment.webappUrl,
    useHeadless: environment.useHeadless,
    useCache: environment.useCache,
  };
}
