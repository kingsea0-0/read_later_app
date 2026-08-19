# Go 后端 API 和数据模型设计

## Question

设计 Hoard Go 后端的 API 架构和数据模型，需要决定以下问题：

1. **项目结构**：Go 项目的目录结构怎么组织？推荐使用什么框架（chi / gin / standard net/http）？
2. **API 端点**：MVP 需要哪些 API 端点？RESTful 设计，包括 bookmarks CRUD、tags 管理、auth 集成等。
3. **认证中间件**：基于 Supabase JWT 的认证中间件如何实现？如何从请求中提取 user_id 并注入 context？
4. **数据层**：pgx 连接池管理、查询构建、事务处理的最佳实践。Go 后端如何实现行级权限过滤（因为 pgx 直连绕过 RLS）？
5. **部署**：API 服务如何部署（Railway/Fly.io）？环境变量管理、健康检查、优雅关停。

## Type

prototype

## Blocked by

- ticket: 001-supabase-research

## Blocks

- ticket: 005-react-frontend-architecture
