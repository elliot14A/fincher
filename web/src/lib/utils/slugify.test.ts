import { describe, expect, it } from 'bun:test'
import { slugify } from './slugify'

describe('slugify', () => {
  it('converts title names into clean kebab-case slugs', () => {
    expect(slugify('Avatar: Fire & Ash')).toBe('avatar-fire-ash')
    expect(slugify('  The Matrix 4  ')).toBe('the-matrix-4')
    expect(slugify('Spider-Man: Across the Spider-Verse')).toBe(
      'spider-man-across-the-spider-verse',
    )
  })

  it('supports prepending custom prefixes', () => {
    expect(slugify('Fire and Ash', 'title')).toBe('title-fire-and-ash')
    expect(slugify('Berlin Synchron', 'vendor')).toBe('vendor-berlin-synchron')
  })

  it('handles empty or punctuation-only strings safely', () => {
    expect(slugify('', 'title')).toBe('title-')
    expect(slugify('---', 'title')).toBe('title-')
    expect(slugify('!!!')).toBe('')
  })
})
