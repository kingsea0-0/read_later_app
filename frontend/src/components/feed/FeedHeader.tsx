import { Bookmark } from 'lucide-react'

type FilterType = 'all' | 'article' | 'video' | 'social' | 'image'

const filterTabs: { key: FilterType; label: string }[] = [
  { key: 'all', label: 'All' },
  { key: 'article', label: 'Articles' },
  { key: 'video', label: 'Videos' },
  { key: 'social', label: 'Social' },
  { key: 'image', label: 'Images' },
]

interface FeedHeaderProps {
  activeFilter: FilterType
  onFilterChange: (filter: FilterType) => void
  totalCount: number
}

export function FeedHeader({ activeFilter, onFilterChange, totalCount }: FeedHeaderProps) {
  return (
    <header className="sticky top-0 z-10 bg-white/95 backdrop-blur-sm border-b border-neutral-200">
      <div className="max-w-6xl mx-auto px-4 sm:px-6">
        <div className="flex items-center justify-between h-14">
          {/* Left: brand */}
          <div className="flex items-center gap-2.5">
            <div className="w-8 h-8 rounded-lg bg-amber-500 flex items-center justify-center">
              <Bookmark size={16} className="text-white" />
            </div>
            <h1 className="text-lg font-bold tracking-tight text-neutral-900">
              <span className="text-amber-500">Hoard</span>
            </h1>
          </div>

          {/* Right: desktop tabs */}
          <nav className="hidden sm:flex items-center gap-1">
            {filterTabs.map((t) => (
              <button
                key={t.key}
                onClick={() => onFilterChange(t.key)}
                className={`px-3 py-1.5 text-sm font-medium rounded-lg transition-colors ${
                  activeFilter === t.key
                    ? 'bg-amber-50 text-amber-600'
                    : 'text-neutral-500 hover:text-neutral-700 hover:bg-neutral-50'
                }`}
              >
                {t.label}
              </button>
            ))}
            <span className="ml-2 text-xs text-neutral-400">{totalCount}</span>
          </nav>
        </div>

        {/* Mobile: scrollable pills */}
        <div className="flex sm:hidden gap-2 pb-3 overflow-x-auto -mx-4 px-4">
          {filterTabs.map((t) => (
            <button
              key={t.key}
              onClick={() => onFilterChange(t.key)}
              className={`px-3 py-1.5 text-xs font-medium rounded-full whitespace-nowrap transition-colors ${
                activeFilter === t.key
                  ? 'bg-amber-500 text-white'
                  : 'bg-neutral-100 text-neutral-500 hover:bg-neutral-200'
              }`}
            >
              {t.label}
            </button>
          ))}
        </div>
      </div>
    </header>
  )
}

export type { FilterType }
