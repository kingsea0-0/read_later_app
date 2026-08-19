# 部署和 DevOps

## Question

确定 Hoard 的部署架构和 DevOps 流程，需要决定以下问题：

1. **域名和 DNS**：选择什么域名（hoard.app / hoard.xxx）？Cloudflare DNS 配置？
2. **Supabase 项目初始化**：创建 Supabase 项目，配置 Auth 提供者，执行数据库迁移脚本。
3. **Go 后端部署**：部署到 Railway / Fly.io，Dockerfile 配置，环境变量管理，健康检查端点。
4. **Vercel 前端部署**：Vercel 项目创建，域名绑定，环境变量，构建配置。
5. **CI/CD**：是否需要 GitHub Actions 自动化部署流程？
6. **监控和日志**：MVP 阶段是否需要监控？最简单的方案是什么？

## Type

task

## Blocked by

- ticket: 004-go-backend-api-design
- ticket: 005-react-frontend-architecture

## Blocks

(none)
