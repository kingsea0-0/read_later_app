// Hoard Extension - Content Script for hoard.app
// Relays OAuth session from Web login to Service Worker

window.addEventListener("message", (event) => {
  if (event.source !== window) return;
  if (event.data && event.data.type === "HOARD_LOGIN") {
    chrome.runtime.sendMessage(event.data);
  }
});

// Also notify service worker that content script is ready
chrome.runtime.sendMessage({ type: "HOARD_CONTENT_READY" });
