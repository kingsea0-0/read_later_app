# Go 后端 API 设计决策（已完成）

## Question

基于 Supabase 调研结论，设计 Hoard 的 Go 后端 API 架构，需要回答以下问题：

1. **项目结构**：Go 项目的目录结构应该如何组织？推荐使用什么框架（chi / gin / 标准库 net/http）？
2. **API 路由**：MVP 需要哪些 API 端点？URL 路径设计、HTTP 方法、请求/响应格式是什么？
3. **认证中间件**：如何实现 JWT 验证中间件？如何从请求中提取 user_id 并传递给 handler？
4. **数据库层**：如何使用 pgx 管理连接池？数据访问层（Repository 模式）如何组织？
5. **部署**：Go 后端如何部署到 Railway/Fly.io？Dockerfile 配置、环境变量管理、健康检查端点。
6. **API 版本化**：URL path 中是否带版本号（/v1/）？

## Type

prototype

## Status

✅ resolved

## Resolution

### 项目结构（按领域分层）

```
backend/
├── cmd/
│   └── server/
│       └── main.go           # 入口，启动 HTTP 服务
├── internal/
│   ├── config/
│   │   └── config.go         # 环境变量配置
│   ├── middleware/
│   │   └── auth.go           # JWT 认证中间件
│   ├── handler/
│   │   ├── auth.go           # 认证相关 handler
│   │   ├── bookmark.go       # 收藏相关 handler
│   │   └── tag.go            # 标签相关 handler
│   ├── repository/
│   │   ├── bookmark.go       # 收藏数据访问
│   │   └── tag.go            # 标签数据访问
│   ├── model/
│   │   ├── bookmark.go       # 收藏数据模型
│   │   └── tag.go            # 标签数据模型
│   └── service/
│       ├── bookmark.go       # 收藏业务逻辑
│       └── tag.go            # 标签业务逻辑
├── migrations/
│   └── 001_init.sql          # 数据库迁移脚本
├── Dockerfile
├── go.mod
└── go.sum
```

### 技术栈

- **HTTP 框架**: gin（用户选择）
- **数据库驱动**: pgx v5（连接池）
- **JWT 验证**: golang-jwt/jwt/v5
- **数据库迁移**: golang-migrate/migrate

### API 端点

所有端点以 `/api/v1/` 为前缀，请求体为 JSON，响应遵循统一格式：

```json
{
  "data": { ... },
  "error": ""
}
```

| 方法 | 路径 | 说明 | 认证 |
|------|------|------|------|
| POST | `/api/v1/auth/session` | 前端登录后，将 Supabase session 传给后端 | 否 |
| DELETE | `/api/v1/auth/session` | 登出 | 是 |
| GET | `/api/v1/bookmarks` | 收藏列表（分页 + 筛选） | 是 |
| POST | `/api/v1/bookmarks` | 创建收藏 | 是 |
| GET | `/api/v1/bookmarks/:id` | 获取单个收藏详情 | 是 |
| PUT | `/api/v1/bookmarks/:id` | 更新收藏（标记已读/加星标等） | 是 |
| DELETE | `/api/v1/bookmarks/:id` | 删除收藏 | 是 |
| GET | `/api/v1/tags` | 获取标签列表 | 是 |
| POST | `/api/v1/tags` | 创建标签 | 是 |
| DELETE | `/api/v1/tags/:id` | 删除标签 | 是 |
| GET | `/healthz` | 健康检查 | 否 |

### 认证中间件

```go
func AuthMiddleware(cfg *config.Config) gin.HandlerFunc {
    return func(c *gin.Context) {
        tokenStr := strings.TrimPrefix(c.GetHeader("Authorization"), "Bearer ")
        if tokenStr == "" {
            c.AbortWithStatusJSON(401, gin.H{"error": "unauthorized"})
            return
        }

        token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
            return []byte(cfg.SupabaseJWTSecret), nil
        })
        if err != nil {
            c.AbortWithStatusJSON(401, gin.H{"error": "invalid token"})
            return
        }

        claims := token.Claims.(jwt.MapClaims)
        userID := claims["sub"].(string)
        c.Set("user_id", userID)
        c.Next()
    }
}
```

### 数据库层（pgx 连接池 + Repository 模式）

```go
// internal/repository/bookmark.go
type BookmarkRepository struct {
    pool *pgxpool.Pool
}

func (r *BookmarkRepository) List(ctx context.Context, userID string, limit, offset int) ([]model.Bookmark, error) {
    rows, err := r.pool.Query(ctx, `
        SELECT id, url, title, site_name, is_archived, is_favorite, created_at
        FROM bookmarks
        WHERE user_id = $1
        ORDER BY created_at DESC
        LIMIT $2 OFFSET $3
    `, userID, limit, offset)
    // ...
}
```

### 部署

- **Dockerfile**: 多阶段构建，最终镜像 ~15MB
- **健康检查**: `GET /healthz` 返回 `{"status": "ok"}`
- **环境变量**: `SUPABASE_URL`, `SUPABASE_JWT_SECRET`, `DATABASE_URL`, `PORT`

## Answer

[详细设计见上方 Resolution 部分]
