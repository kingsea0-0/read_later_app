import { describe, it, expect } from 'vitest'
import { render, screen } from '@testing-library/react'
import { FeedGrid } from './FeedGrid'
import type { Bookmark } from '../../lib/mockData'

function makeBookmark(id: string, overrides: Partial<Bookmark> = {}): Bookmark {
  return {
    id,
    url: `https://example.com/${id}`,
    title: `Title ${id}`,
    description: '',
    favicon_url: '',
    og_image_url: '',
    site_name: 'Example',
    type: 'article',
    is_archived: false,
    is_favorite: false,
    created_at: new Date().toISOString(),
    ...overrides,
  }
}

describe('FeedGrid', () => {
  it('shows an empty state when there are no bookmarks', () => {
    render(<FeedGrid bookmarks={[]} />)
    expect(screen.getByText(/no bookmarks in this category/i)).toBeInTheDocument()
  })

  it('renders a card per bookmark title', () => {
    const bookmarks = [makeBookmark('1'), makeBookmark('2'), makeBookmark('3')]
    render(<FeedGrid bookmarks={bookmarks} />)

    // The first (featured) bookmark renders in both the hero and the mobile
    // grid, so assert on the non-featured titles being present exactly once
    // and the featured title being present.
    expect(screen.getByText('Title 2')).toBeInTheDocument()
    expect(screen.getByText('Title 3')).toBeInTheDocument()
    expect(screen.getAllByText('Title 1').length).toBeGreaterThan(0)
  })

  it('promotes the first bookmark into a hero card', () => {
    const bookmarks = [
      makeBookmark('hero', { title: 'Featured Story' }),
      makeBookmark('2'),
    ]
    render(<FeedGrid bookmarks={bookmarks} />)
    // Featured title is rendered by both the hero (lg) and mobile grid card.
    expect(screen.getAllByText('Featured Story').length).toBe(2)
  })
})
