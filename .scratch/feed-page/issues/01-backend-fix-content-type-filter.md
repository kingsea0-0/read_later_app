# 01 — 后端修复 itoa bug + 增加 content_type 筛选

**What to build:**
- 修复 repository/bookmark.go 中 itoa 函数的 bug（只支持个位数，超过9就出错）
- 在 ListBookmarks API 中增加 content_type 查询参数筛选
- 保证当前构建能通过 go build

**Blocked by:** None — can start immediately

**Status:** ready-for-agent

- [ ] 引入 `strconv` 包，`itoa` 改用 `strconv.Itoa`
- [ ] `BookmarkListParams` 增加 `Type`/`ContentType` 字段绑定查询参数
- [ ] repository `List()` 方法增加 `content_type` WHERE 条件筛选
- [ ] 验证 go build 通过
