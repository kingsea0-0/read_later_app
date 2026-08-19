# 02 — 前端响应式 Feed 页面骨架（杂志风格+PC+手机）

**What to build:**
- 实现响应式 Feed 页面，PC 三列杂志网格，手机单列，平板两列
- 顶部有分类筛选 tabs (All / Articles / Videos)
- 每个卡片展示: OG图片/占位、标题、描述、来源信息、类型标签
- 点击卡片打开原始链接在新标签页
- 无后端/后端连接失败时自动降级使用 mock 数据

**Blocked by:** None — can start immediately

**Status:** ready-for-agent

- [ ] 创建 `components/feed/` 目录，定义组件结构: FeedHeader, FeedCard, FeedGrid
- [ ] 实现响应式布局断点: lg≥1024 → 3列，md 768-1023 → 2列，sm<768 → 1列
- [ ] 实现分类筛选 tabs，PC 用 tabs，手机用 pills
- [ ] 实现卡片组件，展示 OG 图片、标题、描述、来源、类型标签
- [ ] 集成 mock 数据降级逻辑
- [ ] 验证 vite build 通过
