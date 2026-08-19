# Chrome Extension MV3 技术方案调研

## Question


## Status

✅ resolved

## Resolution

Chrome Extension MV3 方案完全可行。
- **架构**：Service Worker 替代 Background Page，Hoard 场景不需要 Content Script 做通用注入，chrome.tabs.query 直接获取 URL 和 Title
- **权限**：只需 activeTab + storage + notifications + alarms，隐私友好
- **一键收藏 UX**：不设 default_popup，点击图标触发 action.onClicked 直接收藏，零弹窗，成功用 notifications 显示短暂通知
- **认证**：Web OAuth 登录后通过 postMessage 传递 session，token 存 chrome.storage.local，alarms 每 50 分钟刷新
- **发布**：$5 注册费，首次审核 1-3 天，最小权限 + 隐私政策 + 无远程代码是审核关键

调研 Chrome Extension Manifest V3 实现方案，需要回答以下问题：

1. **架构**：MV3 下推荐的 extension 架构是什么？Service Worker 如何替代 Background Page？Content Script 和 Service Worker 之间如何通信？
2. **权限**：要读取当前页面 URL 和 Title，需要什么权限？有没有用户隐私方面的限制？
3. **一键收藏流程**：用户点击 extension 图标后，最流畅的 UX 是什么？应该用 Action Popup 还是 Side Panel？如何实现"点击即收藏，无需弹窗确认"？
4. **认证**：Extension 如何与 Supabase Auth 集成？如何在 Extension 中处理 OAuth 登录流程？
5. **API 调用**：Extension 如何调用后端 API？如何安全地存储 token？
6. **Chrome 商店发布**：发布到 Chrome Web Store 的流程和注意事项是什么？

## Type

research

## Blocked by

(none)

## Blocks

- ticket: Go 后端 API 和数据模型设计
- ticket: React 前端信息流架构设计

## Answer

### 1. 架构

#### MV3 整体架构

MV3 强制使用 **Service Worker** 替代 MV2 的 Background Page。Service Worker 是一个事件驱动的、非持久化的 worker，在无事件时会被浏览器休眠以节省内存。

```
┌──────────────────────────────────────────────┐
│                  Browser                      │
│  ┌─────────────────┐    ┌──────────────────┐  │
│  │  Service Worker  │    │  Content Script   │  │
│  │  (background.js) │◄──►│  (per tab)        │  │
│  │  - auth state    │    │  - read DOM       │  │
│  │  - API calls     │    │  - inject UI      │  │
│  │  - chrome.* APIs │    │  - get page info  │  │
│  └────────┬─────────┘    └────────┬─────────┘  │
│           │                       │            │
│           │    messaging API      │            │
│           └───────────────────────┘            │
│  ┌──────────────────────────────────────────┐  │
│  │         Action Popup (popup.html)         │  │
│  │  - ephemeral UI, opens on icon click     │  │
│  └──────────────────────────────────────────┘  │
└──────────────────────────────────────────────┘
```

#### Service Worker 关键特性

- **非持久化**：浏览器随时可以终止 worker。关键状态必须持久化到 `chrome.storage.local`（IndexedDB 也可用但需要额外处理）。
- **唤醒时机**：注册的事件（如 `chrome.action.onClicked`、`alarms`、`webRequest`、`runtime.onMessage`）会唤醒 worker。
- **最长生命周期**：常规事件处理约 30 秒，`chrome.alarms` 最长约 5 分钟。
- **无 DOM 访问**：不能直接操作 DOM，必须通过 Content Script。

#### Content Script ↔ Service Worker 通信

使用 `chrome.runtime.sendMessage` / `chrome.runtime.onMessage`：

```json
// manifest.json — content_scripts 声明
{
  "content_scripts": [{
    "matches": ["<all_urls>"],
    "js": ["content.js"]
  }]
}
```

```javascript
// content.js — 读取当前页面信息并发送给 Service Worker
chrome.runtime.sendMessage({
  type: "PAGE_INFO",
  data: {
    url: window.location.href,
    title: document.title
  }
});
```

```javascript
// background.js (Service Worker) — 接收消息
chrome.runtime.onMessage.addListener((message, sender, sendResponse) => {
  if (message.type === "PAGE_INFO") {
    // 处理页面信息
    sendResponse({ status: "ok" });
  }
  return true; // 保持通道开放以异步响应
});
```

**通信替代方案**：
- `chrome.tabs.query({ active: true, currentWindow: true })` — 从 Service Worker 直接获取当前 tab 信息（包括 `tab.url` 和 `tab.title`），多数情况下不需要 Content Script。
- **一键收藏场景**：使用 `chrome.tabs.query` 即可，无需 Content Script。

### 2. 权限

#### 需要的权限

```json
{
  "permissions": [
    "storage",        // 存储 token 和用户偏好
    "activeTab"       // 获取当前活动标签页的 URL 和 Title
  ],
  "host_permissions": [
    "https://your-api-domain.com/*"  // 调用后端 API
  ]
}
```

#### `activeTab` 权限说明

- **临时权限**：仅在用户点击扩展图标时，扩展获得当前活动标签页的临时访问权限。
- **无需 `tabs` 权限**：`activeTab` 足以读取当前 tab 的 `url` 和 `title`，不需要更宽泛的 `tabs` permission。
- **隐私友好**：`activeTab` 不在安装时提示用户，而是运行时延迟授予，符合 MV3 隐私设计理念。
- MV3 要求 `host_permissions` 显式声明，不再支持从 `permissions` 推断。

#### 用户隐私限制

- Chrome Web Store 对最小权限原则有严格审查。如果只使用 `activeTab` + `storage`，权限提示非常轻量，容易通过审核。
- MV3 取消了远程托管的代码（`eval()`、远程脚本），必须将扩展包内所有代码打包提交。
- 如果扩展在安装后请求额外权限，必须通过 `chrome.permissions.request()` 触发用户确认。

### 3. 一键收藏流程

#### 推荐方案：Action Popup + 无弹窗模式

Hoard 最适合**不使用 Popup**，而是直接监听 `chrome.action.onClicked`：

```json
// manifest.json
{
  "action": {
    "default_icon": {
      "16": "icons/icon16.png",
      "48": "icons/icon48.png",
      "128": "icons/icon128.png"
    }
    // 不设置 default_popup，click 事件会触发 onClicked
  }
}
```

当 `default_popup` 未设置时，点击扩展图标触发 `chrome.action.onClicked`，扩展直接执行收藏操作，不弹出任何 UI。

```javascript
// background.js (Service Worker)
chrome.action.onClicked.addListener(async (tab) => {
  try {
    const token = await getAuthToken();
    if (!token) {
      // 未登录 — 打开登录页面
      await chrome.tabs.create({ url: "https://hoard.app/auth" });
      return;
    }

    const response = await fetch("https://api.hoard.app/v1/bookmarks", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        "Authorization": `Bearer ${token}`
      },
      body: JSON.stringify({
        url: tab.url,
        title: tab.title
      })
    });

    if (response.ok) {
      // 显示成功通知（自动消失，不打扰用户）
      chrome.notifications.create({
        type: "basic",
        iconUrl: "icons/icon128.png",
        title: "Saved to Hoard",
        message: `"${tab.title}" saved`,
        silent: true
      });
    }
  } catch (error) {
    console.error("Failed to save bookmark:", error);
  }
});
```

#### 为什么不用 Side Panel

- Side Panel 需要用户手动打开，不像 Action 点击那么自然。
- Side Panel 占用页面空间，用于"一键收藏"场景太重。
- Side Panel 更适合需要持续交互的场景（如笔记工具、AI 助手）。

#### 为什么不用 Popup

- Popup 需要加载 HTML/CSS/JS，即使只显示半秒也有闪烁。
- "点击即收藏"的最佳 UX 就是零弹窗，用 `onClicked` + `chrome.notifications` 即可。
- 如果需要显示收藏历史或状态，可以用 Action Badge（`chrome.action.setBadgeText`）在图标上显示数字。

#### 完整 UX 流程

```
用户点击图标
  ├─ 未登录 → 打开浏览器新标签页至登录页面
  ├─ 已登录且收藏成功 → 显示短暂通知 + 更新 badge
  └─ 已登录但失败 → 显示错误通知
```

### 4. 认证

#### Supabase Auth 与 Extension 集成方案

Chrome Extension 无法直接使用 Supabase 的 PKCE OAuth 流程（因为 Extension 没有可重定向的 HTTP 端点），推荐以下方案：

##### 方案 A：Web 端 OAuth → 共享 Session（推荐）

用户在 Web 端（`https://hoard.app`）完成 OAuth 登录，Extension 通过 cookie/session 共享或 postMessage 方式获取 token。

```
1. 用户点击 Extension 图标
2. Extension 检测到无 token
3. 打开新标签页到 https://hoard.app/auth?source=extension
4. 用户在 Web 端完成 Google / Apple OAuth 登录
5. Supabase 回传 session 到 Web 前端
6. Web 前端将 session 通过 chrome.runtime.sendMessage 或 URL 参数传递给 Extension
7. Extension 将 token 存入 chrome.storage.local
```

**实现方式**：在 Web 端登录成功后，使用 `chrome.runtime.sendMessage` 或 `chrome.runtime.connect` 将 session 传给 Extension：

```javascript
// Web 端 (在 hoard.app 上)
// 登录成功后，尝试与扩展通信
function sendSessionToExtension(session) {
  // 方式 1: 通过 CustomEvent + DOM 消息（需要 content script 注入）
  window.postMessage({
    source: "hoard-web",
    type: "SUPABASE_SESSION",
    session: {
      access_token: session.access_token,
      refresh_token: session.refresh_token
    }
  }, "*");
}
```

```javascript
// content.js (注入到 hoard.app 的 content script)
window.addEventListener("message", (event) => {
  if (event.data?.source === "hoard-web" && event.data?.type === "SUPABASE_SESSION") {
    chrome.runtime.sendMessage({
      type: "SET_SESSION",
      data: event.data.session
    });
  }
});
```

##### 方案 B：Extension 内嵌 OAuth（授权码流程）

MV3 支持 `chrome.identity.launchWebAuthFlow`，可以在 Extension 内完成 OAuth：

```javascript
// background.js — 使用 chrome.identity API
async function authenticateWithSupabase() {
  const url = `https://your-project.supabase.co/auth/v1/authorize?` +
    `client_id=${SUPABASE_CLIENT_ID}` +
    `&redirect_uri=${chrome.identity.getRedirectURL()}` +
    `&response_type=code` +
    `&scope=openid+profile+email`;

  const redirectUrl = await chrome.identity.launchWebAuthFlow({
    url: url,
    interactive: true
  });

  // 从 redirect URL 提取 authorization code
  const code = new URL(redirectUrl).searchParams.get("code");

  // 交换 code 为 token
  const response = await fetch("https://your-project.supabase.co/auth/v1/token", {
    method: "POST",
    headers: { "Content-Type": "application/x-www-form-urlencoded" },
    body: new URLSearchParams({
      grant_type: "authorization_code",
      code,
      redirect_uri: chrome.identity.getRedirectURL(),
      client_id: SUPABASE_CLIENT_ID
    })
  });

  const tokens = await response.json();
  await chrome.storage.local.set({ supabase_session: tokens });
}
```

**推荐方案 A**，因为：不需要额外配置 OAuth client，Web 端和 Extension 共享一个登录流程，用户体验一致。

### 5. API 调用

#### Token 安全存储

```javascript
// 使用 chrome.storage.local 存储 session
// Key 命名约定：supabase_session
async function saveSession(session) {
  await chrome.storage.local.set({ supabase_session: session });
}

async function getSession() {
  const result = await chrome.storage.local.get("supabase_session");
  return result.supabase_session;
}

async function clearSession() {
  await chrome.storage.local.remove("supabase_session");
}
```

**安全说明**：
- `chrome.storage.local` 数据存储在设备本地，其他扩展无法访问。
- 比 `localStorage` 更安全（`localStorage` 在 Extension 中可用但建议用 `storage` API）。
- 敏感 token 不应放在 `chrome.storage.sync`（会同步到 Google 账号）。
- 不使用 `chrome.storage.session`（只在内存中，Service Worker 唤醒后丢失）。

#### API 调用封装

```javascript
// api.js — 封装后端 API 调用
const API_BASE = "https://api.hoard.app/v1";

async function apiRequest(path, options = {}) {
  const session = await getSession();
  const headers = {
    "Content-Type": "application/json",
    ...options.headers
  };

  if (session?.access_token) {
    headers["Authorization"] = `Bearer ${session.access_token}`;
  }

  const response = await fetch(`${API_BASE}${path}`, {
    ...options,
    headers
  });

  // 如果 401，尝试 refresh token
  if (response.status === 401 && session?.refresh_token) {
    const refreshed = await refreshToken(session.refresh_token);
    if (refreshed) {
      headers["Authorization"] = `Bearer ${refreshed.access_token}`;
      return fetch(`${API_BASE}${path}`, { ...options, headers });
    }
  }

  return response;
}

async function refreshToken(refreshToken) {
  const response = await fetch(
    `${SUPABASE_URL}/auth/v1/token?grant_type=refresh_token`,
    {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        "apikey": SUPABASE_ANON_KEY
      },
      body: JSON.stringify({ refresh_token: refreshToken })
    }
  );

  if (response.ok) {
    const session = await response.json();
    await saveSession(session);
    return session;
  }

  await clearSession();
  return null;
}
```

#### Token 刷新策略

- Supabase access_token 默认有效期 1 小时。
- 设置 `chrome.alarms` 每 50 分钟检查并刷新 token。
- 异常处理：refresh 失败时清除 session，下次点击图标时引导用户重新登录。

```javascript
// 在 Service Worker 启动时注册
chrome.runtime.onInstalled.addListener(() => {
  chrome.alarms.create("refreshToken", { periodInMinutes: 50 });
});

chrome.alarms.onAlarm.addListener(async (alarm) => {
  if (alarm.name === "refreshToken") {
    const session = await getSession();
    if (session?.refresh_token) {
      await refreshToken(session.refresh_token);
    }
  }
});
```

### 6. Chrome 商店发布

#### 发布流程

1. **注册开发者账号**：一次性费用 $5 USD（Google 账号），访问 [Chrome Web Store Developer Dashboard](https://chrome.google.com/webstore/devconsole)。

2. **准备素材**：
   - 图标：128x128 PNG（必需），16x16 和 48x48 可选
   - 至少 1 张截图（1280x800 或 640x400），推荐英文版
   - 小推广图片（440x280）可选
   - 说明文字：详细描述扩展功能

3. **打包扩展**：
   ```bash
   # 将扩展目录打包为 .zip
   zip -r hoard-extension.zip . -x "*.git*" "node_modules/*"
   ```

4. **上传到 CWS**：
   - 在 Dashboard 创建新项目，上传 .zip
   - 填写商店列表信息（名称、描述、类别、语言）
   - 设置隐私政策 URL（必须，Supabase 涉及用户数据）

5. **审核**：
   - 首次审核通常 1-3 个工作日，后续更新更快
   - 审核重点：最小权限原则、隐私政策、代码来源

#### 注意事项

| 类别 | 要求 |
|------|------|
| **权限** | 只请求最小权限。`activeTab` 比 `tabs` 更容易通过审核。`host_permissions` 必须精确。 |
| **隐私政策** | 必须提供。如果使用 Supabase 存储用户数据，政策必须说明什么数据、如何存储、如何删除。 |
| **远程代码** | MV3 禁止远程代码执行。所有代码必须在包内。 |
| **内容安全策略** | 需要使用 CSP 限制 `script-src` 为 `self` 和 `wasm-unsafe-eval`（如果需要）。 |
| **版本号** | 遵循 semver。每次更新必须递增版本号。 |
| **分组发布** | 支持按百分比发布，可用于灰度测试。 |
| **集中式安装** | 如果未来有企业客户，支持通过 `update_url` 托管更新。 |
| **Supabase 特定** | 如果使用 Supabase anon key，它在 Extension 中是可见的，但受 RLS 保护。确保 RLS 策略正确。 |

#### 推荐发布检查清单

- [ ] 权限最小化：使用 `activeTab` 而非 `tabs`
- [ ] 隐私政策文档已编写并托管在 hoard.app/privacy
- [ ] 所有代码在扩展包内，无远程加载
- [ ] 图标和截图已就绪
- [ ] 扩展在多个 Chrome 版本上测试通过
- [ ] 错误处理完善（网络超时、token 过期、API 错误）
- [ ] 使用 `chrome.runtime.setUninstallURL()` 设置卸载后问卷调查页

### 完整 manifest.json 示例

```json
{
  "manifest_version": 3,
  "name": "Hoard — Save for Later",
  "version": "1.0.0",
  "description": "One-click save any page to your Hoard reading list",
  "permissions": [
    "storage",
    "activeTab",
    "notifications",
    "alarms"
  ],
  "host_permissions": [
    "https://api.hoard.app/*",
    "https://*.supabase.co/*"
  ],
  "action": {
    "default_icon": {
      "16": "icons/icon16.png",
      "48": "icons/icon48.png",
      "128": "icons/icon128.png"
    },
    "default_title": "Save to Hoard"
  },
  "background": {
    "service_worker": "background.js",
    "type": "module"
  },
  "content_scripts": [
    {
      "matches": ["https://hoard.app/*"],
      "js": ["content-hoard.js"],
      "run_at": "document_idle"
    }
  ],
  "icons": {
    "16": "icons/icon16.png",
    "48": "icons/icon48.png",
    "128": "icons/icon128.png"
  },
  "minimum_chrome_version": "102"
}
```

### 关键结论

1. **Hoard 不需要 Content Script 做通用页面注入**，只需 `activeTab` 权限通过 `chrome.tabs.query` 获取 URL 和 Title。
2. **Content Script 只需要在 hoard.app 域上注入**，用于接收 Web 端登录后的 session 传递。
3. **使用 `chrome.action.onClicked` 而非 Popup**，实现"零弹窗一键收藏"。
4. **OAuth 走 Web 端登录**，通过 Content Script 传递 session 给 Extension，比在 Extension 内完成 OAuth 更简单可靠。
5. **Token 存 `chrome.storage.local`**，结合 `chrome.alarms` 定时刷新。
6. **发布到 CWS 的重点**：最小权限、隐私政策、无远程代码。
