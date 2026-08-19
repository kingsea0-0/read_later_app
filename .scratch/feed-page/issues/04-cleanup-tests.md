# 04 — 清理原型代码 + 关键路径测试

**What to build:**
- 删除 prototype 目录中的所有原型变体文件（A/B/C/D/E）
- 保留 mock 数据移到 lib 目录
- 对 repository content_type 筛选做简单集成测试
- 对前端 API 层做简单单元测试

**Blocked by:** 03-frontend-infinite-scroll-api

**Status:** blocked

- [ ] 删除 `pages/prototype/` 所有文件
- [ ] 移动 `mockData.ts` 到 `lib/mockData.ts`
- [ ] 清理 FeedPage.tsx，只保留正式实现
- [ ] 添加 repository 集成测试验证 content_type 筛选
- [ ] 添加 API 层单元测试验证参数拼接
- [ ] 全量运行 go test 和 vitest，确认通过
