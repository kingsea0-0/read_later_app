import { describe, it, expect, beforeEach, afterEach, vi } from 'vitest'
import { getTimeAgo } from './time'

describe('getTimeAgo', () => {
  const NOW = new Date('2026-08-19T12:00:00Z').getTime()

  beforeEach(() => {
    vi.useFakeTimers()
    vi.setSystemTime(NOW)
  })

  afterEach(() => {
    vi.useRealTimers()
  })

  function ago(ms: number): string {
    return getTimeAgo(new Date(NOW - ms).toISOString())
  }

  const MIN = 60_000
  const HOUR = 60 * MIN
  const DAY = 24 * HOUR

  it('shows minutes under an hour', () => {
    expect(ago(5 * MIN)).toBe('5m')
    expect(ago(59 * MIN)).toBe('59m')
  })

  it('shows hours under a day', () => {
    expect(ago(1 * HOUR)).toBe('1h')
    expect(ago(23 * HOUR)).toBe('23h')
  })

  it('shows days under a month', () => {
    expect(ago(1 * DAY)).toBe('1d')
    expect(ago(29 * DAY)).toBe('29d')
  })

  it('shows months beyond 30 days', () => {
    expect(ago(30 * DAY)).toBe('1mo')
    expect(ago(90 * DAY)).toBe('3mo')
  })
})
