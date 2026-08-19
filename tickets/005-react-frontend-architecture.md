# React 前端架构设计决策（已完成）

## Question

基于调研结论和 Phase 1 范围（浏览器插件 + Web 阅读端），设计 Hoard 的 React 前端架构，需要回答以下问题：

1. **框架选择**：Vite + React 还是 Next.js SPA？
2. **UI 组件库**：Tailwind CSS + shadcn/ui 还是 MUI？
3. **路由设计**：MVP 需要哪些页面？
4. **状态管理**：是否需要 Zustand/Redux？还是 React Context + hooks 就够用？
5. **组件树**：Card Feed 的核心组件层级设计
6. **认证集成**：前端如何与 Supabase Auth 集成？
7. **API 层**：如何封装对 Go 后端的 API 请求？

## Type

prototype

## Status

✅ resolved

## Resolution

### 技术栈

| 层 | 选择 | 理由 |
|---|---|---|
| 框架 | Vite + React + TypeScript | 轻量、构建快、SPA 天然适配 Vercel |
| UI | Tailwind CSS + shadcn/ui | 灵活可定制，信息流卡片好做 |
| 路由 | React Router v6 | 标准 SPA 路由 |
| 状态管理 | React Context + hooks（MVP） | 够用，后续可迁移到 Zustand |
| API 客户端 | 自定义 fetch 封装 + React Query（可选） | 轻量，自动附带 Bearer token |

### 页面路由

```
/              → 信息流主页（收藏列表）
/auth          → 登录页（Google/Apple 登录按钮）
/auth/callback → OAuth 回调处理页
```

### 核心组件树

```
App
├── AuthProvider (Context)
├── Routes
│   ├── /auth
│   │   └── LoginPage
│   │       ├── GoogleSignInButton
│   │       └── AppleSignInButton
│   ├── /auth/callback
│   │   └── AuthCallbackHandler
│   └── /
│       └── FeedPage
│           ├── FeedHeader (标题 + 筛选 Tab 栏)
│           ├── FeedTabs (全部 / 视频 / 文章)
│           ├── BookmarkCard[]
│           │   ├── CardHeader (favicon + 域名 + 时间)
│           │   ├── CardTitle (可点击跳转原文)
│           │   ├── CardDescription (摘要)
│           │   └── CardActions (标记已读 / 收藏 / 删除)
│           └── InfiniteScroll (无限滚动加载)
```

### 认证集成

```tsx
// lib/supabase.ts
import { createClient } from '@supabase/supabase-js'

export const supabase = createClient(
  import.meta.env.VITE_SUPABASE_URL!,
  import.meta.env.VITE_SUPABASE_ANON_KEY!
)
```

```tsx
// hooks/useAuth.ts
export function useAuth() {
  const [user, setUser] = useState<User | null>(null)
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    supabase.auth.getSession().then(({ data: { session } }) => {
      setUser(session?.user ?? null)
      setLoading(false)
    })
    const { data: { subscription } } = supabase.auth.onAuthStateChange(
      (_event, session) => {
        setUser(session?.user ?? null)
        setLoading(false)
      }
    )
    return () => subscription.unsubscribe()
  }, [])

  const signIn = (provider: 'google' | 'apple') =>
    supabase.auth.signInWithOAuth({ provider,
      options: { redirectTo: `${window.location.origin}/auth/callback` }
    })

  const signOut = () => supabase.auth.signOut()

  return { user, loading, signIn, signOut }
}
```

### API 客户端

```tsx
// lib/api.ts
const API_BASE = import.meta.env.VITE_API_BASE || '/api/v1'

async function apiRequest<T>(path: string, options?: RequestInit): Promise<T> {
  const session = await supabase.auth.getSession()
  const headers: HeadersInit = {
    'Content-Type': 'application/json',
    ...options?.headers,
  }
  if (session.data.session?.access_token) {
    headers['Authorization'] = `Bearer ${session.data.session.access_token}`
  }

  const res = await fetch(`${API_BASE}${path}`, { ...options, headers })
  if (!res.ok) throw new Error(await res.text())
  return res.json()
}
```

## Answer

[详细设计见上方 Resolution 部分]
