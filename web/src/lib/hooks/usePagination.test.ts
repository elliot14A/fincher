import { describe, expect, it } from 'bun:test'
import { DEFAULT_PAGE_LIMIT } from '#/lib/constants'
import { calculatePagination, paginateArray } from './usePagination'

describe('calculatePagination and paginateArray', () => {
  it('calculates total pages and slice bounds accurately for multi-page lists', () => {
    const items = Array.from({ length: 25 }, (_, i) => ({
      id: `item-${i + 1}`,
    }))

    const page1 = paginateArray(items, 1, DEFAULT_PAGE_LIMIT)
    expect(page1.items.length).toBe(10)
    expect(page1.items[0].id).toBe('item-1')
    expect(page1.items[9].id).toBe('item-10')
    expect(page1.calculation.page).toBe(1)
    expect(page1.calculation.totalPages).toBe(3)
    expect(page1.calculation.hasNextPage).toBe(true)
    expect(page1.calculation.hasPrevPage).toBe(false)

    const page2 = paginateArray(items, 2, DEFAULT_PAGE_LIMIT)
    expect(page2.items.length).toBe(10)
    expect(page2.items[0].id).toBe('item-11')
    expect(page2.items[9].id).toBe('item-20')
    expect(page2.calculation.hasNextPage).toBe(true)
    expect(page2.calculation.hasPrevPage).toBe(true)

    const page3 = paginateArray(items, 3, DEFAULT_PAGE_LIMIT)
    expect(page3.items.length).toBe(5)
    expect(page3.items[0].id).toBe('item-21')
    expect(page3.items[4].id).toBe('item-25')
    expect(page3.calculation.hasNextPage).toBe(false)
    expect(page3.calculation.hasPrevPage).toBe(true)
  })

  it('clamps out-of-bounds page requests safely', () => {
    const calcOver = calculatePagination(25, 999, 10)
    expect(calcOver.page).toBe(3)
    expect(calcOver.startIndex).toBe(20)
    expect(calcOver.endIndex).toBe(25)

    const calcUnder = calculatePagination(25, -5, 10)
    expect(calcUnder.page).toBe(1)
    expect(calcUnder.startIndex).toBe(0)
    expect(calcUnder.endIndex).toBe(10)
  })

  it('handles empty arrays gracefully without negative pages or indices', () => {
    const emptyResult = paginateArray([], 1, 10)
    expect(emptyResult.items.length).toBe(0)
    expect(emptyResult.calculation.page).toBe(1)
    expect(emptyResult.calculation.totalPages).toBe(1)
    expect(emptyResult.calculation.hasNextPage).toBe(false)
    expect(emptyResult.calculation.hasPrevPage).toBe(false)
  })
})
