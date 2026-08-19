import { useState, useEffect, useCallback, useRef } from 'react'
import { mockBookmarks } from '../lib/mockData'
import { listBookmarks } from '../lib/api'
import { useInfiniteScroll } from '../lib/useInfiniteScroll'
import { FeedHeader } from '../components/feed/FeedHeader'
import { FeedGrid } from '../components/feed/FeedGrid'
import { Loader2 } from 'lucide-react'
import type { FilterType } from '../components/feed/FeedHeader'
import type { Bookmark as ApiBookmark } from '../lib/api'
import type { Bookmark as MockBookmark } from '../lib/mockData'

const PAGE_SIZE = 20

function mapApiBookmark(b: ApiBookmark): MockBookmark {
  return {
    id: b.id,
    url: b.url,
    title: b.title,
    description: b.description,
    favicon_url: b.favicon_url,
    og_image_url: b.og_image_url,
    site_name: b.site_name,
    type: b.content_type || 'article',
    is_archived: b.is_archived,
    is_favorite: b.is_favorite,
    read_at: b.read_at,
    created_at: b.created_at,
  }
}

function filterMockData(filter: FilterType): MockBookmark[] {
  return filter === 'all'
    ? mockBookmarks
    : mockBookmarks.filter((b) => b.type === filter)
}

export default function FeedPage() {
  const [activeFilter, setActiveFilter] = useState<FilterType>('all')
  const [bookmarks, setBookmarks] = useState<MockBookmark[]>([])
  const [loading, setLoading] = useState(true)
  const [hasMore, setHasMore] = useState(true)
  const [usingMock, setUsingMock] = useState(false)

  // Keep a ref to the current filter so loadMore always reads the latest one
  const filterRef = useRef(activeFilter)
  filterRef.current = activeFilter

  const fetchPage = useCallback(
    async (pageNum: number, filter: FilterType, append: boolean) => {
      setLoading(true)

      try {
        const contentType = filter === 'all' ? undefined : filter
        const res = await listBookmarks(PAGE_SIZE, pageNum * PAGE_SIZE, contentType)

        const mapped = res.bookmarks.map(mapApiBookmark)
        setBookmarks((prev) => (append ? [...prev, ...mapped] : mapped))
        setHasMore(mapped.length === PAGE_SIZE)
        setUsingMock(false)
      } catch (err) {
        console.warn('API unavailable, falling back to mock data:', err)
        const filtered = filterMockData(filter)
        setBookmarks(filtered)
        setHasMore(false)
        setUsingMock(true)
      } finally {
        setLoading(false)
      }
    },
    []
  )

  // Reset on filter change
  useEffect(() => {
    setBookmarks([])
    setHasMore(true)
    setUsingMock(false)
    setLoading(true)
    fetchPage(0, activeFilter, false)
  }, [activeFilter, fetchPage])

  const loadMore = useCallback(() => {
    // Don't load more while loading, when there's no more data, or when using mock data
    if (loading || !hasMore || usingMock) return
    // We need the current page number — derive it from bookmark count
    const nextPage = Math.floor(bookmarks.length / PAGE_SIZE)
    fetchPage(nextPage, filterRef.current, true)
  }, [loading, hasMore, usingMock, bookmarks.length, fetchPage])

  const sentinelRef = useInfiniteScroll({
    onLoadMore: loadMore,
    hasMore: hasMore && !usingMock,
  })

  return (
    <div className="min-h-screen bg-neutral-50">
      <FeedHeader
        activeFilter={activeFilter}
        onFilterChange={setActiveFilter}
        totalCount={bookmarks.length}
      />
      <FeedGrid bookmarks={bookmarks} />

      {/* Sentinel for infinite scroll + loading spinner */}
      {!usingMock && (
        <div ref={sentinelRef} className="flex justify-center py-6">
          {loading && (
            <Loader2 size={20} className="animate-spin text-neutral-400" />
          )}
          {!hasMore && bookmarks.length > 0 && (
            <span className="text-xs text-neutral-400">All bookmarks loaded</span>
          )}
        </div>
      )}
    </div>
  )
}
