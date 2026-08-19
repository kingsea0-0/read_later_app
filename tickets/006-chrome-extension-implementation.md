# Chrome Extension 实现方案决策（已完成）

## Question

基于 Chrome Extension 调研结论，设计 Hoard 的 Chrome Extension 实现方案，需要回答以下问题：

1. **项目结构**：Chrome Extension 的目录结构和构建工具怎么选？
2. **manifest.json**：完整的 manifest 配置
3. **Service Worker**：background.js 的完整实现
4. **Content Script**：hoard.app 域上的 content script
5. **构建流程**：如何打包为 .zip
6. **图标设计**：需要哪些尺寸的图标

## Type

prototype

## Status

✅ resolved

## Resolution

### 项目结构

```
extension/
├── manifest.json
├── background.js        # Service Worker
├── content-hoard.js     # 注入 hoard.app 的 Content Script
├── icons/
│   ├── icon16.png
│   ├── icon48.png
│   └── icon128.png
└── build.sh             # 打包脚本
```

纯手写，不需要构建工具。代码量很小，TypeScript 会增加复杂度。

### manifest.json

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
  }
}
```

### Service Worker 核心逻辑

```javascript
// background.js
const API_BASE = 'https://api.hoard.app/v1'

// 点击图标一键收藏
chrome.action.onClicked.addListener(async (tab) => {
  const session = await getSession()
  if (!session?.access_token) {
    await chrome.tabs.create({ url: 'https://hoard.app/auth' })
    return
  }

  try {
    await apiRequest('/bookmarks', {
      method: 'POST',
      body: JSON.stringify({ url: tab.url, title: tab.title })
    })
    chrome.notifications.create({
      type: 'basic',
      iconUrl: 'icons/icon128.png',
      title: 'Saved to Hoard',
      message: `"${tab.title}" saved`,
      silent: true
    })
  } catch (e) {
    console.error('Failed to save:', e)
  }
})

// 定时刷新 token（每 50 分钟）
chrome.runtime.onInstalled.addListener(() => {
  chrome.alarms.create('refreshToken', { periodInMinutes: 50 })
})

chrome.alarms.onAlarm.addListener(async (alarm) => {
  if (alarm.name === 'refreshToken') {
    const session = await getSession()
    if (session?.refresh_token) {
      await refreshToken(session.refresh_token)
    }
  }
})

// 从 Web 端接收 session
chrome.runtime.onMessage.addListener((message, sender, sendResponse) => {
  if (message.type === 'SET_SESSION') {
    chrome.storage.local.set({ supabase_session: message.data })
    sendResponse({ status: 'ok' })
  }
})
```

### Content Script（hoard.app 域）

```javascript
// content-hoard.js
window.addEventListener('message', (event) => {
  if (event.data?.source === 'hoard-web' && event.data?.type === 'SUPABASE_SESSION') {
    chrome.runtime.sendMessage({
      type: 'SET_SESSION',
      data: event.data.session
    })
  }
})
```

### 构建流程

```bash
# 打包为 .zip，用于 Chrome Web Store 上传
cd extension
zip -r hoard-extension.zip . -x "*.git*"
```

## Answer

[详细设计见上方 Resolution 部分]
