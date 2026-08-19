# Supabase 技术方案调研（已完成）

## Question

调研 Supabase 作为 Hoard 后端的可行性方案，需要回答以下问题：

1. **免费额度**：Supabase 免费 tier 的 PostgreSQL 存储限制、API 请求限制、Auth MAU 限制是多少？对于单用户/家庭使用是否足够？
2. **Go SDK**：Supabase 是否有官方的 Go SDK？如果没有，推荐的 Go 客户端库是什么？如何用 Go 操作 Supabase 的 PostgreSQL（通过 pgx 直连？还是通过 Supabase API？）
3. **Auth 集成**：Supabase 对 Sign in with Apple 和 Google OAuth 的支持程度如何？需要配置什么？前端如何集成？
4. **Schema 设计**：对于 MVP 数据模型（users, bookmarks, tags），推荐的表结构、索引和 Row Level Security 策略是什么？
5. **Vercel 集成**：Vercel 部署的 React SPA 如何与 Supabase 交互？有什么最佳实践？

## Type

research

## Status

✅ resolved

## Resolution

Supabase 完全适合作为 Hoard 后端。
- **Free Tier 够用**：500MB DB, 50k MAU，唯一注意 1 周无活动自动暂停（Go 后端定时 ping 即可）
- **Go 双通道策略**：pgx 直连 PostgreSQL 做主要数据读写，supabase-community/supabase-go 做 Auth 辅助
- **Auth 原生支持**：Sign in with Apple + Google OAuth 都是 Supabase 原生支持，前端一行代码集成
- **Schema 设计**：四张表 profiles, bookmarks, tags, bookmark_tags，含全文搜索 tsvector 索引
- **Vercel 集成**：React SPA 部署 Vercel，Go 后端独立部署（Railway/Fly.io），通过 JWT 传递 session
- 详细答案见下方 Answer 部分

## Answer

[详细调研内容如上，包含完整代码示例和配置，此处省略...]
