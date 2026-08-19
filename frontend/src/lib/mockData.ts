export interface Bookmark {
  id: string
  url: string
  title: string
  description?: string
  favicon_url?: string
  og_image_url?: string
  site_name?: string
  type: 'article' | 'video' | 'social' | 'image'
  is_archived: boolean
  is_favorite: boolean
  read_at?: string | null
  created_at: string
}

const now = new Date()
const daysAgo = (d: number) => {
  const t = new Date(now)
  t.setDate(t.getDate() - d)
  return t.toISOString()
}

export const mockBookmarks: Bookmark[] = [
  {
    id: '1',
    url: 'https://www.youtube.com/watch?v=example1',
    title: 'How to Build a Read-Later App in 2026 — Full Tutorial',
    description: 'A comprehensive walkthrough of building a modern read-later application with Go, React, and Chrome Extensions. Covers architecture, auth, and deployment.',
    favicon_url: 'https://www.youtube.com/favicon.ico',
    og_image_url: 'https://images.unsplash.com/photo-1555066931-4365d14bab8c?w=800',
    site_name: 'YouTube',
    type: 'video',
    is_archived: false,
    is_favorite: true,
    created_at: daysAgo(0),
  },
  {
    id: '2',
    url: 'https://zhuanlan.zhihu.com/p/example2',
    title: '系统设计面试：如何设计一个类似 Instapaper 的阅读服务',
    description: '本文详细分析了设计一个可扩展的稍后阅读服务所需考虑的核心组件，包括数据存储、内容抓取、全文搜索和推荐系统。',
    favicon_url: 'https://zhuanlan.zhihu.com/favicon.ico',
    site_name: '知乎专栏',
    type: 'article',
    is_archived: false,
    is_favorite: false,
    created_at: daysAgo(1),
  },
  {
    id: '3',
    url: 'https://twitter.com/example/status/3',
    title: 'Thread: Why I switched from bookmarks to a read-later workflow',
    description: 'A deep dive into how I process 200+ links per week using a simple read-later system. Tips for reducing information overload.',
    favicon_url: 'https://twitter.com/favicon.ico',
    og_image_url: 'https://images.unsplash.com/photo-1611162617474-5b21e879e113?w=800',
    site_name: 'Twitter',
    type: 'social',
    is_archived: false,
    is_favorite: true,
    created_at: daysAgo(2),
  },
  {
    id: '4',
    url: 'https://www.bilibili.com/video/example4',
    title: '【硬核】一行代码都不用写，用 AI 搭建个人知识库',
    description: '保姆级教程，手把手教你用最新的 AI 工具搭建属于自己的知识管理系统。',
    favicon_url: 'https://www.bilibili.com/favicon.ico',
    og_image_url: 'https://images.unsplash.com/photo-1677442136019-21780ecad995?w=800',
    site_name: 'Bilibili',
    type: 'video',
    is_archived: false,
    is_favorite: false,
    created_at: daysAgo(3),
  },
  {
    id: '5',
    url: 'https://www.ft.com/content/example5',
    title: 'The Art of Deep Reading in the Age of Information Overload',
    description: 'How the most productive people in tech manage their reading habits. A look at the tools and techniques that enable deep focus.',
    favicon_url: 'https://www.ft.com/favicon.ico',
    og_image_url: 'https://images.unsplash.com/photo-1455390582262-044cdead277a?w=800',
    site_name: 'Financial Times',
    type: 'article',
    is_archived: false,
    is_favorite: false,
    created_at: daysAgo(4),
  },
  {
    id: '6',
    url: 'https://feishu.cn/doc/example6',
    title: '团队知识管理最佳实践：从 Evernote 到飞书文档的迁移之路',
    description: '分享我们团队在知识管理工具选型上的四次迭代经验，以及最终选择飞书文档的原因和使用技巧。',
    favicon_url: 'https://feishu.cn/favicon.ico',
    site_name: '飞书文档',
    type: 'article',
    is_archived: false,
    is_favorite: false,
    created_at: daysAgo(5),
  },
  {
    id: '7',
    url: 'https://www.xiaohongshu.com/explore/example7',
    title: '我的桌面布置 2.0｜极简主义书桌分享',
    description: '终于找到了最适合自己的桌面布置方案，给大家分享一下好物清单和布局思路～',
    favicon_url: 'https://www.xiaohongshu.com/favicon.ico',
    og_image_url: 'https://images.unsplash.com/photo-1593062096033-9a26b09da705?w=800',
    site_name: '小红书',
    type: 'social',
    is_archived: false,
    is_favorite: false,
    created_at: daysAgo(6),
  },
  {
    id: '8',
    url: 'https://weibo.com/example8',
    title: '推荐！这 10 个开源项目让你的开发效率翻倍 🔥',
    description: '整理了我最近发现的好用开源项目，涵盖 CLI 工具、前端框架和数据库工具。',
    favicon_url: 'https://weibo.com/favicon.ico',
    site_name: '微博',
    type: 'social',
    is_archived: false,
    is_favorite: false,
    created_at: daysAgo(7),
  },
  {
    id: '9',
    url: 'https://www.architect.io/posts/example9',
    title: 'Building a Serverless Read-Later API with Cloudflare Workers and D1',
    description: 'Learn how to build a cost-effective, globally distributed API for saving and retrieving bookmarks using Cloudflare\'s edge network.',
    favicon_url: 'https://www.architect.io/favicon.ico',
    og_image_url: 'https://images.unsplash.com/photo-1558494949-ef010cbdcc31?w=800',
    site_name: 'Architect.io',
    type: 'article',
    is_archived: false,
    is_favorite: false,
    created_at: daysAgo(8),
  },
  {
    id: '10',
    url: 'https://www.youtube.com/watch?v=example10',
    title: 'I Tried 10 Read-Later Apps So You Don\'t Have To',
    description: 'From Pocket to Matter to Omnivore — here\'s what each app does well and where they fall short.',
    favicon_url: 'https://www.youtube.com/favicon.ico',
    og_image_url: 'https://images.unsplash.com/photo-1504711434969-e33886168d6c?w=800',
    site_name: 'YouTube',
    type: 'video',
    is_archived: false,
    is_favorite: false,
    created_at: daysAgo(10),
  },
  {
    id: '11',
    url: 'https://blog.example.com/11',
    title: 'Why Your Bookmark Collection Is a Graveyard (And How to Fix It)',
    description: 'Most people save links and never look at them again. Here\'s the simple system that finally made my bookmark collection useful.',
    favicon_url: 'https://blog.example.com/favicon.ico',
    site_name: 'Productivity Blog',
    type: 'article',
    is_archived: false,
    is_favorite: false,
    created_at: daysAgo(14),
  },
  {
    id: '12',
    url: 'https://www.bilibili.com/video/example12',
    title: '【Vlog】我的信息处理工作流：从输入到输出',
    description: '分享一下我每天如何处理海量信息，包括阅读、笔记、写作的完整流程。',
    favicon_url: 'https://www.bilibili.com/favicon.ico',
    og_image_url: 'https://images.unsplash.com/photo-1499750310107-5fef28a66643?w=800',
    site_name: 'Bilibili',
    type: 'video',
    is_archived: false,
    is_favorite: false,
    created_at: daysAgo(21),
  },
]
