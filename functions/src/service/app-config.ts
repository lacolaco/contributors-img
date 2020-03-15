export interface ApplicationConfig {
  webappUrl: string;
  useHeadless: boolean;
  useCache: boolean;
}

export function getApplicationConfig(): ApplicationConfig {
  return {
    webappUrl: process.env.WEBAPP_URL || 'https://contributors-img.firebaseapp.com',
    useHeadless: process.env.USE_HEADLESS ? process.env.USE_HEADLESS === 'true' : true,
    useCache: process.env.USE_CACHE ? process.env.USE_CACHE === 'true' : true,
  };
}
