import { ExternalLink, Play, BookOpen, MessageCircle, Image } from 'lucide-react'
import type { Bookmark } from '../../lib/mockData'

function getTimeAgo(dateStr: string): string {
  const now = Date.now()
  const then = new Date(dateStr).getTime()
  const diff = now - then
  const mins = Math.floor(diff / 60000)
  if (mins < 60) return `${mins}m`
  const hours = Math.floor(diff / 60)
  if (hours < 24) return `${hours}h`
  const days = Math.floor(diff / 24)
  if (days < 30) return `${days}d`
  return `${Math.floor(days / 30)}mo`
}

const typeConfig = {
  video: { icon: Play, label: 'Video', color: 'text-red-500', bg: 'bg-red-50' },
  article: { icon: BookOpen, label: 'Article', color: 'text-blue-500', bg: 'bg-blue-50' },
  social: { icon: MessageCircle, label: 'Social', color: 'text-purple-500', bg: 'bg-purple-50' },
  image: { icon: Image, label: 'Image', color: 'text-emerald-600', bg: 'bg-emerald-50' },
}

export function FeedCard({ bookmark }: { bookmark: Bookmark }) {
  const domain = bookmark.url
    ? new URL(bookmark.url).hostname.replace('www.', '')
    : ''
  const timeAgo = getTimeAgo(bookmark.created_at)
  const type = typeConfig[bookmark.type]
  const TypeIcon = type.icon

  return (
    <article className="group bg-white rounded-xl border border-neutral-200 overflow-hidden hover:shadow-md hover:border-neutral-300 transition-all duration-200">
      {/* Image or placeholder */}
      <a
        href={bookmark.url}
        target="_blank"
        rel="noopener noreferrer"
        className="block"
      >
        {bookmark.og_image_url ? (
          <div className="aspect-[16/9] bg-neutral-100 overflow-hidden">
            <img
              src={bookmark.og_image_url}
              alt=""
              className="w-full h-full object-cover group-hover:scale-105 transition-transform duration-300"
              loading="lazy"
            />
          </div>
        ) : (
          <div className="aspect-[16/9] bg-gradient-to-br from-amber-100 to-amber-50 flex items-center justify-center">
            <span className="text-4xl font-bold text-amber-300 select-none">
              {bookmark.title.charAt(0)}
            </span>
          </div>
        )}
      </a>

      <div className="p-4">
        {/* Type badge */}
        <div className="flex items-center gap-2 mb-2">
          <span className={`inline-flex items-center gap-1 text-[10px] font-semibold uppercase tracking-wider ${type.color}`}>
            <TypeIcon size={10} />
            {type.label}
          </span>
          {bookmark.is_favorite && (
            <span className="text-amber-400 text-xs">★</span>
          )}
        </div>

        {/* Title — clicking opens the link */}
        <a
          href={bookmark.url}
          target="_blank"
          rel="noopener noreferrer"
          className="text-sm font-semibold text-neutral-900 group-hover:text-amber-600 line-clamp-2 transition-colors"
        >
          {bookmark.title}
        </a>

        {/* Description — max 2 lines */}
        {bookmark.description && (
          <p className="mt-1.5 text-xs text-neutral-500 line-clamp-2 leading-relaxed">
            {bookmark.description}
          </p>
        )}

        {/* Source row: favicon + domain + time */}
        <div className="flex items-center gap-2 mt-3 pt-3 border-t border-neutral-100">
          <img
            src={bookmark.favicon_url || `https://www.google.com/s2/favicons?domain=${domain}`}
            alt=""
            className="w-3.5 h-3.5 rounded"
            onError={(e) => {
              (e.target as HTMLImageElement).style.display = 'none'
            }}
          />
          <span className="text-xs font-medium text-neutral-500 truncate">
            {bookmark.site_name || domain}
          </span>
          <span className="text-neutral-300">·</span>
          <span className="text-xs text-neutral-400">{timeAgo}</span>
          <a
            href={bookmark.url}
            target="_blank"
            rel="noopener noreferrer"
            className="ml-auto text-neutral-400 hover:text-neutral-600 transition-colors"
            title="Open original"
          >
            <ExternalLink size={12} />
          </a>
        </div>
      </div>
    </article>
  )
}
