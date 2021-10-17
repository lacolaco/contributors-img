export function createPager(pageLimit: number) {
  return (size: number): number[] => {
    const pages: number[] = [];
    let i = size;
    while (i > 0) {
      const page = Math.min(i, pageLimit);
      pages.push(page);
      i -= page;
    }
    return pages;
  };
}
