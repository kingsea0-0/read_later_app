import { describe, it, expect, vi } from 'vitest'
import { render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { FeedHeader } from './FeedHeader'

describe('FeedHeader', () => {
  it('renders all filter tabs', () => {
    render(<FeedHeader activeFilter="all" onFilterChange={() => {}} totalCount={0} />)
    // Labels appear twice (desktop tabs + mobile pills); getAllByText tolerates that.
    expect(screen.getAllByText('All').length).toBeGreaterThan(0)
    expect(screen.getAllByText('Articles').length).toBeGreaterThan(0)
    expect(screen.getAllByText('Videos').length).toBeGreaterThan(0)
    expect(screen.getAllByText('Social').length).toBeGreaterThan(0)
    expect(screen.getAllByText('Images').length).toBeGreaterThan(0)
  })

  it('fires onFilterChange with the correct key when a tab is clicked', async () => {
    const onFilterChange = vi.fn()
    render(<FeedHeader activeFilter="all" onFilterChange={onFilterChange} totalCount={0} />)

    // Click the first "Videos" button (desktop tab).
    await userEvent.click(screen.getAllByText('Videos')[0])
    expect(onFilterChange).toHaveBeenCalledWith('video')
  })

  it('fires onFilterChange with "article" for the Articles tab', async () => {
    const onFilterChange = vi.fn()
    render(<FeedHeader activeFilter="all" onFilterChange={onFilterChange} totalCount={0} />)

    await userEvent.click(screen.getAllByText('Articles')[0])
    expect(onFilterChange).toHaveBeenCalledWith('article')
  })

  it('shows the total count', () => {
    render(<FeedHeader activeFilter="all" onFilterChange={() => {}} totalCount={42} />)
    expect(screen.getByText('42')).toBeInTheDocument()
  })
})
