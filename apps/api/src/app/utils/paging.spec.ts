import { createPager } from './paging';

describe('createPager', () => {
  const pager = createPager(100);

  test('0 becomes []', () => {
    const pages = pager(0);
    expect(pages).toEqual([]);
  });

  test('1 becomes [1]', () => {
    const pages = pager(1);
    expect(pages).toEqual([1]);
  });

  test('100 becomes [100]', () => {
    const pages = pager(100);
    expect(pages).toEqual([100]);
  });

  test('101 becomes [100, 1]', () => {
    const pages = pager(101);
    expect(pages).toEqual([100, 1]);
  });

  test('101 becomes [100, 1]', () => {
    const pages = pager(101);
    expect(pages).toEqual([100, 1]);
  });

  test('201 becomes [100, 100, 1]', () => {
    const pages = pager(201);
    expect(pages).toEqual([100, 100, 1]);
  });
});
