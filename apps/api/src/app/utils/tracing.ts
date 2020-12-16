import { get } from '@google-cloud/trace-agent';

export function addTracingLabels(labels: Record<string, string>) {
  const rootSpan = get().getCurrentRootSpan();
  for (const [key, value] of Object.entries(labels)) {
    rootSpan.addLabel(key, value);
  }
}

export function runWithTracing<T>(name: string, fn: () => Promise<T>): Promise<T> {
  const span = get().createChildSpan({ name });
  return fn().finally(() => {
    span.endSpan();
  });
}
