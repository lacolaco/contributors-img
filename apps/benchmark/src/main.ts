import { Contributor } from '@lib/core';
import { createMockContributor } from '@lib/core/testing';
import { createJsRenderer, createRustRenderer } from '@lib/renderer';

const jsRenderer = createJsRenderer();
const rustRenderer = createRustRenderer();

async function renderer_js(contributors: Contributor[]) {
  await jsRenderer.render(contributors);
}

async function renderer_rust(contributors: Contributor[]) {
  await rustRenderer.render(contributors);
}

async function benchmark_js(times: number, contributors: Contributor[]) {
  performance.mark('js:start');
  for (let index = 0; index < times; index++) {
    await renderer_js(contributors);
  }
  performance.mark('js:stop');
  performance.measure('js', 'js:start', 'js:stop');
}

async function benchmark_rust(times: number, contributors: Contributor[]) {
  performance.mark('rust:start');
  for (let index = 0; index < times; index++) {
    await renderer_rust(contributors);
  }
  performance.mark('rust:stop');
  performance.measure('rust', 'rust:start', 'rust:stop');
}

async function main() {
  console.log('start benchmark');
  const times = 1000;

  async function run(dataSize: number) {
    console.log(`=== data size: ${dataSize} ===`);
    const contributors = new Array(dataSize).fill(null).map((_, i) => createMockContributor({ login: `login${i}` }));
    await benchmark_js(times, contributors);
    await benchmark_rust(times, contributors);
    for (const measure of performance.getEntriesByType('measure')) {
      console.log(`${measure.name}: ${(measure.duration / times).toPrecision(2)}ms`);
    }
    performance.clearMarks();
    performance.clearMeasures();
  }

  await run(1);
  await run(2);
  await run(5);
  await run(10);
  await run(20);
  await run(50);
  await run(100);
}

main();
