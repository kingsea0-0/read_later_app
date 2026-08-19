import { describe, it, expect } from 'vitest'
import { render, screen } from '@testing-library/react'
import { FeedCard } from './FeedCard'
import type { Bookmark } from '../../lib/mockData'

function makeBookmark(overrides: Partial<Bookmark> = {}): Bookmark {
  return {
    id: '1',
    url: 'https://example.com/post',
    title: 'Test Title',
    description: 'A short description',
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

describe('FeedCard', () => {
  it('renders the article type badge for article content', () => {
    render(<FeedCard bookmark={makeBookmark({ type: 'article' })} />)
    expect(screen.getByText('Article')).toBeInTheDocument()
  })

  it('renders the video type badge for video content', () => {
    render(<FeedCard bookmark={makeBookmark({ type: 'video' })} />)
    expect(screen.getByText('Video')).toBeInTheDocument()
  })

  it('renders the social type badge for social content', () => {
    render(<FeedCard bookmark={makeBookmark({ type: 'social' })} />)
    expect(screen.getByText('Social')).toBeInTheDocument()
  })

  it('renders the image type badge for image content', () => {
    render(<FeedCard bookmark={makeBookmark({ type: 'image' })} />)
    expect(screen.getByText('Image')).toBeInTheDocument()
  })

  it('links the title to the original url in a new tab', () => {
    render(<FeedCard bookmark={makeBookmark({ url: 'https://example.com/post' })} />)
    const titleLink = screen.getByRole('link', { name: 'Test Title' })
    expect(titleLink).toHaveAttribute('href', 'https://example.com/post')
    expect(titleLink).toHaveAttribute('target', '_blank')
    expect(titleLink).toHaveAttribute('rel', 'noopener noreferrer')
  })

  it('shows a favorite star only when is_favorite is true', () => {
    const { rerender } = render(<FeedCard bookmark={makeBookmark({ is_favorite: false })} />)
    expect(screen.queryByText('★')).not.toBeInTheDocument()

    rerender(<FeedCard bookmark={makeBookmark({ is_favorite: true })} />)
    expect(screen.getByText('★')).toBeInTheDocument()
  })

  it('falls back to a title initial placeholder when there is no og image', () => {
    render(<FeedCard bookmark={makeBookmark({ title: 'Zebra', og_image_url: '' })} />)
    expect(screen.getByText('Z')).toBeInTheDocument()
  })

  it('renders the description when provided', () => {
    render(<FeedCard bookmark={makeBookmark({ description: 'Hello world desc' })} />)
    expect(screen.getByText('Hello world desc')).toBeInTheDocument()
  })
})
