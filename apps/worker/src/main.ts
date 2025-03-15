import { startServer } from './internal/server';

async function main() {
  try {
    const port = process.env.PORT ? parseInt(process.env.PORT) : 3333;
    const appEnv = process.env.APP_ENV || 'development';
    await startServer(port, appEnv);
  } catch (error) {
    console.error('Failed to start server:', error);
    process.exit(1);
  }
}

main();
