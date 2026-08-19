// Hoard Extension - Service Worker

const API_BASE = "http://localhost:8080/api/v1";
const TOKEN_REFRESH_INTERVAL = 50; // minutes

// ========== Token Management ==========

async function getToken() {
  const result = await chrome.storage.local.get(["accessToken", "refreshToken"]);
  return result.accessToken || null;
}

async function setToken(accessToken, refreshToken) {
  await chrome.storage.local.set({ accessToken, refreshToken });
}

async function clearToken() {
  await chrome.storage.local.remove(["accessToken", "refreshToken"]);
}

async function refreshToken() {
  const result = await chrome.storage.local.get(["refreshToken"]);
  if (!result.refreshToken) return false;

  try {
    const res = await fetch(`${API_BASE}/auth/refresh`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ refresh_token: result.refreshToken }),
    });
    if (!res.ok) throw new Error("refresh failed");
    const data = await res.json();
    await setToken(data.access_token, data.refresh_token);
    return true;
  } catch {
    await clearToken();
    return false;
  }
}

// ========== API Request ==========

async function apiRequest(path, options = {}) {
  const token = await getToken();
  const headers = { ...options.headers };
  if (token) {
    headers["Authorization"] = `Bearer ${token}`;
  }
  if (!(options.body instanceof FormData)) {
    headers["Content-Type"] = "application/json";
  }

  const res = await fetch(`${API_BASE}${path}`, { ...options, headers });

  if (res.status === 401 && token) {
    const refreshed = await refreshToken();
    if (refreshed) {
      const newToken = await getToken();
      headers["Authorization"] = `Bearer ${newToken}`;
      const retryRes = await fetch(`${API_BASE}${path}`, { ...options, headers });
      return retryRes;
    }
    // Redirect to login
    await clearToken();
    await chrome.tabs.create({ url: "http://localhost:5173/auth" });
    throw new Error("session expired");
  }

  return res;
}

// ========== Action: Save Current Tab ==========

chrome.action.onClicked.addListener(async (tab) => {
  const url = tab.url;
  const title = tab.title || "Untitled";

  // Skip chrome://, about:, etc.
  if (!url.startsWith("http")) return;

  try {
    const res = await apiRequest("/bookmarks", {
      method: "POST",
      body: JSON.stringify({ url, title }),
    });

    if (res.ok) {
      await chrome.notifications.create({
        type: "basic",
        iconUrl: "icons/icon48.png",
        title: "Saved to Hoard",
        message: `"${title.substring(0, 80)}"`,
        priority: 1,
      });
    } else {
      const err = await res.json();
      await chrome.notifications.create({
        type: "basic",
        iconUrl: "icons/icon48.png",
        title: "Save Failed",
        message: err.error || "Unknown error",
        priority: 2,
      });
    }
  } catch (err) {
    await chrome.notifications.create({
      type: "basic",
      iconUrl: "icons/icon48.png",
      title: "Save Failed",
      message: "Connection error. Check your network.",
      priority: 2,
    });
  }
});

// ========== Token Refresh Alarm ==========

chrome.alarms.create("tokenRefresh", { periodInMinutes: TOKEN_REFRESH_INTERVAL });

chrome.alarms.onAlarm.addListener(async (alarm) => {
  if (alarm.name === "tokenRefresh") {
    const token = await getToken();
    if (token) {
      await refreshToken();
    }
  }
});

// ========== Web Login Session Relay ==========

chrome.runtime.onMessageExternal.addListener(async (message, sender, sendResponse) => {
  if (message.type === "HOARD_LOGIN" && message.accessToken) {
    await setToken(message.accessToken, message.refreshToken || "");
    sendResponse({ success: true });
  }
});
