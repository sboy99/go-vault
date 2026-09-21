import { describe, expect, it } from 'vitest'
import {
  formatBytes,
  formatDuration,
  formatRelativeTime,
  truncateHash,
} from '@/lib/utils/format'

describe('formatBytes', () => {
  it('formats zero and small values', () => {
    expect(formatBytes(0)).toBe('0 B')
    expect(formatBytes(512)).toBe('512 B')
  })

  it('formats larger units', () => {
    expect(formatBytes(1024)).toBe('1.00 KB')
    expect(formatBytes(1536)).toBe('1.50 KB')
    expect(formatBytes(1048576)).toBe('1.00 MB')
  })

  it('handles invalid input', () => {
    expect(formatBytes(-1)).toBe('—')
    expect(formatBytes(Number.NaN)).toBe('—')
  })
})

describe('formatRelativeTime', () => {
  it('returns dash for empty', () => {
    expect(formatRelativeTime(null)).toBe('—')
    expect(formatRelativeTime(undefined)).toBe('—')
  })

  it('formats recent times', () => {
    const now = Date.parse('2026-01-01T12:00:00Z')
    expect(formatRelativeTime('2026-01-01T11:59:30Z', now)).toMatch(/second/)
  })
})

describe('formatDuration', () => {
  it('formats durations', () => {
    expect(formatDuration('2026-01-01T00:00:00Z', '2026-01-01T00:00:00.500Z')).toBe('500ms')
    expect(formatDuration('2026-01-01T00:00:00Z', '2026-01-01T00:00:05Z')).toBe('5s')
    expect(formatDuration('2026-01-01T00:00:00Z', '2026-01-01T00:02:05Z')).toBe('2m 5s')
  })
})

describe('truncateHash', () => {
  it('truncates long hashes', () => {
    expect(truncateHash('abcdef0123456789', 8)).toBe('abcdef01…')
    expect(truncateHash('short')).toBe('short')
    expect(truncateHash(undefined)).toBe('—')
  })
})
