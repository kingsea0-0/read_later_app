# 03 — 前端无限滚动 + API 对接

**What to build:**
- API 层增加 `listBookmarks` 的 content_type 参数支持
- Feed 页面实现无限滚动分页加载（下拉到底加载下一页）
- 切换筛选条件时重置分页重新加载
- 展示加载状态和 loading spinner

**Blocked by:** 01-backend-fix-content-type-filter, 02-frontend-responsive-feed

**Status:** blocked

- [ ] 更新 `lib/api.ts` 中 `listBookmarks` 支持 contentType 参数
- [ ] 实现无限滚动 hook 或使用 react-infinite-scroll-hook
- [ ] 整合分页状态到 Feed 页面
- [ ] 切换筛选重置分页
- [ ] 添加加载状态展示
