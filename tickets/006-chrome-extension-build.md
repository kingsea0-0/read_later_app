# Chrome Extension 构建

## Question

基于 Chrome Extension 调研结果，构建 Hoard 的浏览器扩展，需要决定以下问题：

1. **项目结构**：Chrome Extension 的目录结构怎么组织？是否需要构建工具（如 Vite/webpack）？
2. **manifest.json**：MVP 版本的 manifest 配置，权限声明，host_permissions。
3. **Service Worker**：核心逻辑实现：点击图标收藏、token 管理、API 调用、定时刷新 token。
4. **Web 端通信**：hoard.app 域上的 Content Script，用于接收 Web 登录后的 session 传递。
5. **图标和素材**：扩展图标、通知图标的准备。
6. **打包和发布**：构建脚本、打包命令、Chrome Web Store 发布清单。

## Type

prototype

## Blocked by

- ticket: 002-chrome-extension-research
- ticket: 004-go-backend-api-design

## Blocks

(none)
