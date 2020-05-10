import { environment } from '../../environments/environment.prod';

export interface ApplicationConfig {
  webappUrl: string;
  useHeadless: boolean;
  useCache: boolean;
}

export function getApplicationConfig(): ApplicationConfig {
  return {
    webappUrl: environment.WEBAPP_URL,
    useHeadless: environment.USE_HEADLESS,
    useCache: environment.USE_CACHE,
  };
}
