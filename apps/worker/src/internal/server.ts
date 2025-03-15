import { serve } from '@hono/node-server';
import { Hono } from 'hono';
import { queryFeaturedRepositories, queryUsageStats } from './query';
import { saveFeaturedRepositories, saveUsageStats } from './store';

export async function startServer(port: number, appEnv: string): Promise<void> {
  const app = new Hono();

  app.all('/update-featured-repositories', async (c) => {
    try {
      const featuredRepositories = await queryFeaturedRepositories();
      const usageStats = await queryUsageStats();
      const updatedAt = new Date();

      await saveFeaturedRepositories(appEnv, featuredRepositories, updatedAt);
      await saveUsageStats(appEnv, usageStats, updatedAt);

      return c.json({ message: 'Data updated successfully' }, 200);
    } catch (error) {
      console.error('Failed to update data:', error);
      return c.json({ error: 'Failed to update data' }, 500);
    }
  });

  app.all('*', (c) => c.json({ error: 'Not Found' }, 404));

  serve({ fetch: app.fetch, port }, () => {
    console.log(`Server running at http://127.0.0.1:${port}/`);
  });
}
