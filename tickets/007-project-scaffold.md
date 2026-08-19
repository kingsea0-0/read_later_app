# 项目脚手架搭建（已完成）

## Question

Hoard 项目需要初始化以下代码仓库和目录结构，确保后续开发可以顺利进行。

## Type

task

## Status

✅ completed

## Blocked by

- [Go 后端 API 设计决策](tickets/004-go-api-design.md) ✅
- [React 前端架构设计决策](tickets/005-react-frontend-architecture.md) ✅
- [Chrome Extension 实现方案决策](tickets/006-chrome-extension-implementation.md) ✅

## Resolution

### 目录结构

```
read_later_app/
├── backend/          # Go 后端 (gin + pgx)
├── frontend/         # React 前端 (Vite + Tailwind + shadcn)
├── extension/        # Chrome Extension (MV3)
├── ios/              # iOS 原生壳 (Phase 2)
├── migrations/       # 数据库迁移脚本
├── tickets/          # 决策 ticket
├── .env.example      # 环境变量模板
└── .gitignore        # Git 忽略规则
```

### 后端 (backend/)

- **Go 1.24** + gin v1.10 + pgx v5 + golang-jwt v5
- 目录：cmd/server, internal/{config,handler,middleware,model,repository,service}
- **main.go**: 完整入口，含 graceful shutdown、路由注册、健康检查
- **handler/auth.go**: 认证 session 管理
- **handler/bookmark.go**: 收藏 CRUD
- **handler/tag.go**: 标签 CRUD
- **middleware/auth.go**: JWT Bearer token 验证中间件
- **model/bookmark.go**: Bookmark 数据模型 + 请求/响应结构体
- **model/tag.go**: Tag 数据模型
- **repository/bookmark.go**: 收藏数据访问层（含动态更新 SQL）
- **repository/tag.go**: 标签数据访问层
- **service/bookmark.go**: 收藏业务逻辑
- **service/tag.go**: 标签业务逻辑
- **config/config.go**: 环境变量配置加载
- **Dockerfile**: 多阶段构建 ~15MB
- ✅ `go build ./cmd/server` 通过

### 前端 (frontend/)

- Vite 8 + React 19 + TypeScript + Tailwind CSS
- 已安装：react-router-dom, @supabase/supabase-js, lucide-react
- **App.tsx**: 路由配置（/, /auth, /auth/callback）
- **pages/FeedPage.tsx**: 信息流主页（含无限滚动、Tab 筛选）
- **pages/LoginPage.tsx**: 登录页（Google/Apple 按钮）
- **pages/AuthCallbackPage.tsx**: OAuth 回调处理页
- **contexts/AuthContext.tsx**: 认证上下文（Supabase Auth）
- **lib/supabase.ts**: Supabase 客户端初始化
- **lib/api.ts**: Go 后端 API 封装
- ✅ `vite build` 通过

### Chrome Extension (extension/)

- MV3 纯手写，无需构建工具
- **manifest.json**: 最小权限声明（storage, activeTab, notifications, alarms）
- **background.js**: Service Worker — 点击图标收藏、token 管理、定时刷新
- **content-hoard.js**: Web 端 session 传递
- **icons/**: 16/48/128 PNG 图标（深色底 + 琥珀色 H 字母）
- **build.sh**: 打包脚本

### 数据库迁移 (migrations/)

- **001_init.sql**: 完整 schema（profiles, bookmarks, tags, bookmark_tags）
- 含全文搜索 tsvector、RLS 策略、触发器

### 环境变量

- `.env.example`: 前端 + 后端环境变量模板

### 已验证

- ✅ Go 后端构建通过
- ✅ 前端 build 通过
- ✅ Chrome Extension 完整文件结构就绪
- ✅ 数据库迁移脚本就绪
