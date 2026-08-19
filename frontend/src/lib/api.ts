import { supabase } from './supabase'

const API_BASE = import.meta.env.VITE_API_URL || 'http://localhost:8080/api/v1'

export interface Bookmark {
  id: string
  user_id: string
  url: string
  title: string
  description?: string
  favicon_url?: string
  og_image_url?: string
  site_name?: string
  content_type: 'article' | 'video' | 'social' | 'image'
  is_archived: boolean
  is_favorite: boolean
  read_at?: string | null
  created_at: string
  updated_at: string
}

export interface BookmarkListResponse {
  bookmarks: Bookmark[]
  total: number
}

export interface Tag {
  id: string
  user_id: string
  name: string
  color?: string
  created_at: string
}

async function getAccessToken(): Promise<string | null> {
  const { data } = await supabase.auth.getSession()
  return data.session?.access_token ?? null
}

async function request<T>(path: string, options: RequestInit = {}): Promise<T> {
  const token = await getAccessToken()
  const headers: Record<string, string> = {
    'Content-Type': 'application/json',
    ...(options.headers as Record<string, string> || {}),
  }
  if (token) {
    headers['Authorization'] = `Bearer ${token}`
  }

  const res = await fetch(`${API_BASE}${path}`, { ...options, headers })

  if (!res.ok) {
    if (res.status === 204) {
      return undefined as T
    }
    const err = await res.json().catch(() => ({ error: res.statusText }))
    throw new Error(err.error || 'Request failed')
  }

  if (res.status === 204) {
    return undefined as T
  }

  return res.json()
}

export async function listBookmarks(limit = 20, offset = 0, contentType?: string): Promise<BookmarkListResponse> {
  let path = `/bookmarks?limit=${limit}&offset=${offset}`
  if (contentType) {
    path += `&type=${encodeURIComponent(contentType)}`
  }
  return request<BookmarkListResponse>(path)
}

export async function createBookmark(url: string, title: string, contentType?: string): Promise<Bookmark> {
  const body: Record<string, string> = { url, title }
  if (contentType) {
    body.content_type = contentType
  }
  const res = await request<{ data: Bookmark }>('/bookmarks', {
    method: 'POST',
    body: JSON.stringify(body),
  })
  return res.data
}

export async function updateBookmark(id: string, updates: Partial<Bookmark>): Promise<Bookmark> {
  const res = await request<{ data: Bookmark }>(`/bookmarks/${id}`, {
    method: 'PUT',
    body: JSON.stringify(updates),
  })
  return res.data
}

export async function deleteBookmark(id: string): Promise<void> {
  await request(`/bookmarks/${id}`, { method: 'DELETE' })
}

export async function listTags(): Promise<{ tags: Tag[]; total: number }> {
  return request<{ tags: Tag[]; total: number }>('/tags')
}

export async function createTag(name: string, color?: string): Promise<Tag> {
  const res = await request<{ data: Tag }>('/tags', {
    method: 'POST',
    body: JSON.stringify({ name, color }),
  })
  return res.data
}

export async function deleteTag(id: string): Promise<void> {
  await request(`/tags/${id}`, { method: 'DELETE' })
}
