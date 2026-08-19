# Feed 页面 + 后端 API 端点实现

## Problem Statement

用户通过浏览器插件或手机分享将链接保存到 Hoard 后，需要一个浏览这些收藏内容的阅读界面。目前只有原型代码（5 个变体），没有正式的生产级 Feed 页面，也没有一个可用的后端 API 端点来提供数据。用户无法在浏览器中看到自己的收藏，也无法验证端到端的"保存→查看"流程。

## Solution

实现一个响应式 Feed 页面（PC 杂志风格 + 手机单列布局），连接后端 `GET /api/v1/bookmarks` 端点，支持分页加载和按内容类型筛选。页面在本地上无后端时可用 mock 数据工作，连上后端后自动切换到真实数据。

## User Stories

1. 作为 Hoard 用户，我打开首页就能看到我的收藏列表，以便快速浏览最近保存的内容
2. 作为 Hoard 用户，我可以在 PC 上看到三列杂志风格的卡片网格，以便高效浏览大量内容
3. 作为 Hoard 用户，我可以在手机上看到单列全宽卡片，以便在窄屏上舒适地阅读
4. 作为 Hoard 用户，我可以通过顶部分类栏筛选「All / Articles / Videos / Social」，以便只看某一类内容
5. 作为 Hoard 用户，我点击卡片可以直接跳转到原始链接，以便阅读完整内容
6. 作为 Hoard 用户，我可以看到卡片上的来源网站图标、标题、描述和保存时间，以便快速判断内容价值
7. 作为 Hoard 用户，我下拉页面可以加载更多收藏（无限滚动或分页），以便浏览所有内容
8. 作为 Hoard 用户，我在没有后端时也能看到 mock 示例数据，以便本地开发调试
9. 作为后端开发者，我可以通过 `GET /api/v1/bookmarks?type=video&limit=20&offset=0` 获取分页的收藏列表，以便前端消费
10. 作为后端开发者，JWT 认证中间件可以正确解析 token 并提取 user_id，以便保护用户数据
11. 作为后端开发者，API 响应格式统一为 `{ bookmarks: [...], total: N }`，以便前端易于解析
12. 作为用户，已收藏的卡片显示星标状态，以便我区分重要内容

## Implementation Decisions

### 1. 后端：Bookmark List API 增强

**修改文件**：
- `backend/internal/model/bookmark.go` — `BookmarkListParams` 增加 `ContentType` 字段，`BookmarkListResponse` 统一格式
- `backend/internal/repository/bookmark.go` — `List()` 方法增加 `content_type` 筛选参数
- `backend/internal/handler/bookmark.go` — `ListBookmarks` handler 传递 `ContentType` 到 repository
- `backend/internal/handler/auth.go` — 完善 JWT 中间件，解析 Supabase JWT 并注入 user_id

**修复 Bug**：
- `repository/bookmark.go` 中的 `itoa()` 函数仅支持个位数，改用 `strconv.Itoa`

**API 响应格式统一**：
- `List`: `{ bookmarks: [...], total: N }` ✅ 已有
- `Create/Get/Update`: 当前返回 `{ data: bookmark }` — 保持现状，前端通过 `res.data` 解包
- `Delete`: 返回 `204 No Content` 而不是 `{ data: "deleted" }`

**content_type 筛选**：
- `content_type` 可选值：`article`, `video`, `social`, `image`
- 空值或不传时返回所有类型
- 对应数据库 `bookmarks.content_type` 字段

### 2. 后端：cmd/main.go 路由注册

创建 `backend/cmd/main.go`，注册所有路由：
- `GET /api/v1/bookmarks` — ListBookmarks
- `POST /api/v1/bookmarks` — CreateBookmark
- `GET /api/v1/bookmarks/:id` — GetBookmark
- `PUT /api/v1/bookmarks/:id` — UpdateBookmark
- `DELETE /api/v1/bookmarks/:id` — DeleteBookmark
- `GET /api/v1/tags` — ListTags
- `POST /api/v1/tags` — CreateTag
- `DELETE /api/v1/tags/:id` — DeleteTag

### 3. 前端：响应式 Feed 页面（原型合并）

**设计决策（来自原型）**：
- 采用的布局：原型 C（PC 杂志风格）+ 原型 D（手机单列）
- 通过 Tailwind 响应式断点合并：
  - `sm` (< 768px)：单列，全宽卡片
  - `md` (768px - 1023px)：两列网格
  - `lg` (≥ 1024px)：三列网格 + 顶部 Hero 大图

**组件结构**：
```
pages/FeedPage.tsx          — 主页面，负责数据加载
components/feed/
  FeedHeader.tsx            — 顶部导航栏 + 分类筛选 tabs
  HeroCard.tsx              — 首条内容的 Hero 大图（PC 端）
  FeedCard.tsx              — 通用卡片组件
  FeedGrid.tsx              — 响应式网格容器
  FilterBar.tsx             — 分类 pills（mobile）/ tabs（PC）
```

**数据流**：
- 尝试调用 `listBookmarks()` API
- 如果 API 调用失败（无后端），优雅降级到 mock 数据
- 滚动到底部自动加载下一页（infinite scroll）
- 切换分类筛选时重置分页

**样式**：
- Tailwind CSS，使用 `@apply` 提取公共样式
- 采用从原型继承的琥珀色主题（`text-amber-500`）
- 卡片圆角 `rounded-xl`，阴影 `hover:shadow-md`
- 图片使用 `aspect-[16/9]` 和 `object-cover`
- 类型标签用颜色区分：视频红色、文章蓝色、社交紫色

### 4. 前端：API 层增强

**`frontend/src/lib/api.ts`** 现有 API 层增强：
- `listBookmarks()` 支持 `contentType` 参数
- 添加 `listBookmarksInfinite()` 支持分页

### 5. 原型清理

完成实现后：
- 删除 `frontend/src/pages/prototype/` 目录所有文件
- 删除所有 Variant 组件
- 删除 PrototypeSwitcher
- 将 mock 数据移到 `frontend/src/lib/mockData.ts`

## Testing Decisions

### 测试策略（MVP）

MVP 阶段采用最轻量有效的测试覆盖：

**后端**：
- Repository 层：集成测试，使用测试 PostgreSQL 数据库（docker-compose 或 Supabase local）
  - 测试 `List()` 带 `content_type` 筛选
  - 测试 `Create()` 和 `GetByID()` 的 CRUD 流程
  - 测试 `Update()` 和 `Delete()` 的行级权限过滤
- Handler 层：HTTP 测试，使用 `net/http/httptest` + gin 测试模式
  - 测试认证头缺失时返回 401
  - 测试正常请求的响应格式

**前端**：
- 组件渲染测试：使用 Vitest + @testing-library/react
  - 测试 FeedCard 渲染不同 content_type 时的正确样式
  - 测试 FeedGrid 在不同 screen size 下的列数
  - 测试 FilterBar 切换时正确触发回调
- API 层测试：mock fetch，测试请求参数和响应解析

### 现有 seams

项目中尚无现有测试，这是第一个测试覆盖。优先测试关键的业务逻辑路径。

## Out of Scope

- AI 摘要功能 — Phase 3
- 用户自定义分组/文件夹 — Phase 2
- 标签系统 UI — 仅后端 API 存在，前端不展示
- 搜索功能 — 后续 spec
- 二级详情页 — 点击卡片直接跳转原文
- iOS Share Extension — Phase 2
- Chrome Extension 的前端交互 — 独立的 spec
- 部署配置 — 在 007 脚手架中已规划

## Further Notes

- 原型变体 A、B、E 在本次实现中不采用，但保留在 git 历史中供参考
- 后端 API 当前没有 content_type 自动检测逻辑（由 Chrome Extension 或 iOS Share Extension 在创建时传入）
- 用户认证前置条件：用户需要先通过 `/auth` 页面登录，获取 Supabase JWT
- cmd/main.go 需要从环境变量读取 DATABASE_URL 和 SUPABASE_JWT_SECRET
- itoa bug 修复：`strconv.Itoa(i)` 替代 `string(rune('0'+i))`
