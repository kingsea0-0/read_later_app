import { Play, BookOpen } from 'lucide-react'
import { FeedCard } from './FeedCard'
import type { Bookmark } from '../../lib/mockData'
import { getTimeAgo } from '../../lib/time'

interface FeedGridProps {
  bookmarks: Bookmark[]
}

export function FeedGrid({ bookmarks }: FeedGridProps) {
  if (bookmarks.length === 0) {
    return (
      <div className="text-center py-16 text-neutral-400 text-sm">
        No bookmarks in this category yet.
      </div>
    )
  }

  const featured = bookmarks[0]
  const rest = bookmarks.slice(1)

  return (
    <div className="max-w-6xl mx-auto px-4 sm:px-6 py-6">
      {/* Hero — only on lg screens */}
      {featured && (
        <div className="hidden lg:block mb-8">
          <HeroCard bookmark={featured} />
        </div>
      )}

      {/* Grid: lg=3 col, md=2 col, sm=1 col */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-5">
        {/* On mobile/tablet, include featured in the grid as the first card */}
        {featured && (
          <div className="lg:hidden">
            <FeedCard bookmark={featured} />
          </div>
        )}
        {rest.map((bookmark) => (
          <FeedCard key={bookmark.id} bookmark={bookmark} />
        ))}
      </div>
    </div>
  )
}

function HeroCard({ bookmark }: { bookmark: Bookmark }) {
  const domain = bookmark.url
    ? new URL(bookmark.url).hostname.replace('www.', '')
    : ''
  const timeAgo = getTimeAgo(bookmark.created_at)

  return (
    <div className="relative bg-neutral-900 rounded-2xl overflow-hidden">
      {/* Background image */}
      {bookmark.og_image_url && (
        <div className="absolute inset-0">
          <img
            src={bookmark.og_image_url}
            alt=""
            className="w-full h-full object-cover opacity-40"
          />
          <div className="absolute inset-0 bg-gradient-to-t from-neutral-900 via-neutral-900/60 to-transparent" />
        </div>
      )}

      <div className="relative p-8 md:p-12">
        {/* Type badge */}
        <span className="inline-flex items-center gap-1.5 px-2.5 py-1 text-xs font-medium bg-white/10 backdrop-blur-sm text-white rounded-full mb-4">
          {bookmark.type === 'video' ? <Play size={12} /> : <BookOpen size={12} />}
          {bookmark.type}
        </span>

        <h2 className="text-2xl md:text-3xl font-bold text-white leading-tight max-w-2xl">
          <a
            href={bookmark.url}
            target="_blank"
            rel="noopener noreferrer"
            className="hover:text-amber-300 transition-colors"
          >
            {bookmark.title}
          </a>
        </h2>

        {bookmark.description && (
          <p className="mt-3 text-sm text-white/70 max-w-xl leading-relaxed line-clamp-2">
            {bookmark.description}
          </p>
        )}

        {/* Meta row */}
        <div className="flex items-center gap-3 mt-5">
          <img
            src={bookmark.favicon_url || `https://www.google.com/s2/favicons?domain=${domain}`}
            alt=""
            className="w-5 h-5 rounded"
            onError={(e) => {
              (e.target as HTMLImageElement).style.display = 'none'
            }}
          />
          <span className="text-sm font-medium text-white/90">
            {bookmark.site_name || domain}
          </span>
          <span className="text-white/40">·</span>
          <span className="text-sm text-white/60">{timeAgo}</span>
          <a
            href={bookmark.url}
            target="_blank"
            rel="noopener noreferrer"
            className="ml-auto inline-flex items-center gap-1.5 px-4 py-2 text-sm font-medium text-white bg-amber-500 hover:bg-amber-600 rounded-xl transition-colors"
          >
            Open
          </a>
        </div>
      </div>
    </div>
  )
}
