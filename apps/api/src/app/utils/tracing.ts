import { get } from '@google-cloud/trace-agent';

export function runWithTracing<T>(name: string, fn: () => Promise<T>): Promise<T> {
  const span = get().createChildSpan({ name });
  return fn().finally(() => {
    span.endSpan();
  });
}
