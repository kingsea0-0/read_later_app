import { describe, it, expect, vi, beforeEach } from 'vitest'

// Mock supabase so api.ts doesn't throw when calling getAccessToken
vi.mock('./supabase', () => ({
  supabase: {
    auth: {
      getSession: vi.fn().mockResolvedValue({ data: { session: null } }),
    },
  },
}))

describe('listBookmarks', () => {
  beforeEach(() => {
    // Create a fresh mock for each test
    const mockFetch = vi.fn().mockResolvedValue({
      ok: true,
      json: async () => ({ bookmarks: [], total: 0 }),
    })
    vi.stubGlobal('fetch', mockFetch)
  })

  it('calls /bookmarks with limit and offset when no contentType', async () => {
    const { listBookmarks } = await import('./api')
    const result = await listBookmarks(20, 0)

    const fetchMock = vi.mocked(fetch)
    expect(fetchMock).toHaveBeenCalledTimes(1)
    const url = fetchMock.mock.calls[0][0] as string
    expect(url).toContain('/bookmarks?limit=20&offset=0')
    expect(url).not.toContain('&type=')
    expect(result).toEqual({ bookmarks: [], total: 0 })
  })

  it('appends type=video when contentType is video', async () => {
    const { listBookmarks } = await import('./api')
    await listBookmarks(20, 0, 'video')

    const fetchMock = vi.mocked(fetch)
    const url = fetchMock.mock.calls[0][0] as string
    expect(url).toContain('/bookmarks?limit=20&offset=0')
    expect(url).toContain('&type=video')
  })

  it('appends type=article when contentType is article', async () => {
    const { listBookmarks } = await import('./api')
    await listBookmarks(20, 0, 'article')

    const url = vi.mocked(fetch).mock.calls[0][0] as string
    expect(url).toContain('&type=article')
  })

  it('encodes special characters in contentType', async () => {
    const { listBookmarks } = await import('./api')
    await listBookmarks(20, 0, 'video/podcast')

    const url = vi.mocked(fetch).mock.calls[0][0] as string
    expect(url).toContain('&type=video%2Fpodcast')
  })
})
