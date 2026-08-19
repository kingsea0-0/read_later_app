# iOS Share Extension 技术方案调研

## Question


## Status

✅ resolved

## Resolution

iOS Share Extension 方案可行，工作量很小。
- **最小实现**：3 个文件 ~80 行 Swift，声明 NSExtensionActivationSupportsWebPageWithMaxCount 接收 URL + Title
- **API 调用**：推荐 Share Extension 先写入 App Group 共享目录的 JSON 队列，主 App 在 applicationDidBecomeActive 时同步到后端，避免 30 秒超时限制
- **认证共享**：App Groups（文件共享）+ Keychain Sharing（token 共享），kSecAttrAccessibleAfterFirstUnlock 确保解锁后可用
- **WKWebView Shell**：2 个文件 ~100 行 Swift，全屏 WKWebView 加载 SPA，拦截 hoard:// URL Scheme 实现 JS↔Native 通信
- **App Store**：需 $99 开发者账号，主 App 和 Extension 用同一 Bundle ID 体系

调研 iOS Share Extension 的最小实现方案，需要回答以下问题：

1. **Share Extension 基础**：iOS Share Extension 如何接收来自其他 App 的分享（URL + Title）？最小的实现代码量是多少？
2. **API 调用**：Share Extension 中如何调用后端 API 保存数据？如何处理网络请求和认证？
3. **认证流程**：用户在 iOS App 中完成登录后，Share Extension 如何共享登录态？App Groups 和 Keychain Sharing 的配置方式？
4. **WKWebView Shell**：如何实现一个极薄的 iOS App，内部用 WKWebView 加载 React SPA？URL Scheme 和 Deep Link 如何处理？
5. **App Store 发布**：发布到 App Store 的要求（Apple Developer Program 账号、证书配置等）

## Type

research

## Blocked by

(none)

## Blocks

- ticket: Go 后端 API 和数据模型设计

## Answer

### 1. Share Extension 基础

**工作原理**：iOS Share Extension 是一个独立的 target，通过 `NSExtensionActivationRule` 声明自己能接收哪些数据类型。当用户在系统分享菜单中选中你的 Extension 时，iOS 将分享内容（如 `NSItemProvider`）传入 Extension 进程。

**最小文件清单**：

| 文件 | 作用 |
|------|------|
| `ShareExtension/Info.plist` | 声明 Extension 类型和激活规则 |
| `ShareExtension/ShareViewController.swift` | 主入口，处理分享内容 |
| `ShareExtension/Base.lproj/MainInterface.storyboard` | 可选，可用纯代码替代 |

**Info.plist 关键配置**：

```xml
<key>NSExtension</key>
<dict>
    <key>NSExtensionAttributes</key>
    <dict>
        <key>NSExtensionActivationRule</key>
        <dict>
            <key>NSExtensionActivationSupportsWebURLWithMaxCount</key>
            <integer>1</integer>
            <key>NSExtensionActivationSupportsWebPageWithMaxCount</key>
            <integer>1</integer>
        </dict>
    </dict>
    <key>NSExtensionMainStoryboard</key>
    <string>MainInterface</string>
    <key>NSExtensionPointIdentifier</key>
    <string>com.apple.share-services</string>
</dict>
```

`NSExtensionActivationSupportsWebPageWithMaxCount: 1` 表示接收 Safari 等浏览器分享的网页（含 URL + Title），`NSExtensionActivationSupportsWebURLWithMaxCount: 1` 表示接收纯 URL 分享。两者同时设置可覆盖更多场景。

**ShareViewController.swift 最小实现**：

```swift
import UIKit
import Social

class ShareViewController: SLComposeServiceViewController {

    override func isContentValid() -> Bool {
        return true
    }

    override func didSelectPost() {
        guard let item = extensionContext?.inputItems.first as? NSExtensionItem,
              let attachments = item.attachments else {
            cancel()
            return
        }

        let group = DispatchGroup()
        var urlString: String?
        var title: String?

        for provider in attachments {
            if provider.hasItemConformingToTypeIdentifier("public.url") {
                group.enter()
                provider.loadItem(forTypeIdentifier: "public.url", options: nil) { (item, error) in
                    if let url = item as? URL {
                        urlString = url.absoluteString
                    }
                    group.leave()
                }
            }
            if provider.hasItemConformingToTypeIdentifier("public.plain-text") {
                group.enter()
                provider.loadItem(forTypeIdentifier: "public.plain-text", options: nil) { (item, error) in
                    if let text = item as? String {
                        title = text
                    }
                    group.leave()
                }
            }
        }

        group.notify(queue: .main) {
            self.saveToHoard(url: urlString, title: title)
        }
    }

    private func saveToHoard(url: String?, title: String?) {
        // 调用 API 或写入本地队列
        extensionContext?.completeRequest(returningItems: nil, completionHandler: nil)
    }

    override func configurationItems() -> [Any]! {
        return []
    }
}
```

> 注意：`SLComposeServiceViewController` 会自动显示系统样式的弹出窗口，包含标题输入框，用户可在分享前编辑标题。如需完全自定义 UI，可继承 `UIViewController`。

**Xcode 配置步骤**：
1. File → New → Target → iOS → Share Extension
2. 配置 Product Name（如 `ShareExtension`）
3. 自动生成上述文件结构

---

### 2. API 调用

**网络请求方式**：Share Extension 中使用标准的 `URLSession` 发起网络请求，与主 App 完全相同。

```swift
struct APIClient {
    static let baseURL = "https://api.hoard.app"

    static func saveBookmark(url: String, title: String, token: String, completion: @escaping (Bool) -> Void) {
        var request = URLRequest(url: URL(string: "\(baseURL)/api/bookmarks")!)
        request.httpMethod = "POST"
        request.setValue("application/json", forHTTPHeaderField: "Content-Type")
        request.setValue("Bearer \(token)", forHTTPHeaderField: "Authorization")

        let body: [String: Any] = [
            "url": url,
            "title": title
        ]
        request.httpBody = try? JSONSerialization.data(withJSONObject: body)

        URLSession.shared.dataTask(with: request) { data, response, error in
            let success = (response as? HTTPURLResponse)?.statusCode == 200 || (response as? HTTPURLResponse)?.statusCode == 201
            completion(success)
        }.resume()
    }
}
```

**关键约束**：

- **时间限制**：Share Extension 必须在约 30 秒内完成操作，否则系统会超时终止。
- **后台模式不可用**：Share Extension 不支持 `background URLSession`，所有请求必须在 Extension 存活期内完成。
- **无 UI 阻塞**：如果网络慢，考虑在 `didSelectPost` 中立即返回，但 Extension 可能随时被系统终止。

**推荐方案：先存本地再同步**：

对于 Hoard 的用例，更稳妥的做法是：
1. Share Extension 收到分享后，先写入 **App Group 共享目录**（本地 JSON 或 Core Data）
2. 主 App 启动时或通过 `CFNotificationCenter` 触发同步，调用后端 API 提交
3. 这种方式不受 Extension 超时限制，用户体验更好

```swift
func enqueueForSync(url: String, title: String) {
    let sharedDir = FileManager.default
        .containerURL(forSecurityApplicationGroupIdentifier: "group.com.hoard.app")!
        .appendingPathComponent("pending_bookmarks.json")

    var pending = (try? JSONSerialization.jsonObject(with: Data(contentsOf: sharedDir)) as? [[String: String]]) ?? []
    pending.append(["url": url, "title": title, "createdAt": ISO8601DateFormatter().string(from: Date())])

    if let data = try? JSONSerialization.data(withJSONObject: pending) {
        try? data.write(to: sharedDir)
    }
}
```

主 App 在 `AppDelegate.applicationDidBecomeActive` 读取该文件并同步到后端。

---

### 3. 认证流程

**问题**：iOS App 和 Share Extension 运行在不同进程中，无法直接共享内存中的登录态。需要通过 **App Groups** + **Keychain Sharing** 解决。

#### App Groups 配置

1. 在 Apple Developer Portal 或 Xcode Capabilities 中启用 App Groups
2. 设置 App Group ID，如 `group.com.hoard.app`
3. 主 App target 和 Share Extension target 都勾选该 App Group

```
Target → Signing & Capabilities → + → App Groups
→ 勾选 "group.com.hoard.app"
```

#### Keychain Sharing 配置

同等重要：在 Keychain 中共享 auth token。

```
Target → Signing & Capabilities → + → Keychain Sharing
→ 设置 Keychain Group 如 "com.hoard.app"
```

**代码实现**：

```swift
import Security

struct KeychainService {
    static let serviceName = "com.hoard.app.auth"
    static let keychainGroup = "com.hoard.app"

    static func save(token: String) {
        guard let data = token.data(using: .utf8) else { return }
        let query: [String: Any] = [
            kSecClass as String: kSecClassGenericPassword,
            kSecAttrService as String: serviceName,
            kSecAttrAccessGroup as String: keychainGroup,
            kSecAttrAccount as String: "auth_token",
            kSecValueData as String: data,
            kSecAttrAccessible as String: kSecAttrAccessibleAfterFirstUnlock
        ]
        SecItemDelete(query as CFDictionary)
        SecItemAdd(query as CFDictionary, nil)
    }

    static func getToken() -> String? {
        let query: [String: Any] = [
            kSecClass as String: kSecClassGenericPassword,
            kSecAttrService as String: serviceName,
            kSecAttrAccessGroup as String: keychainGroup,
            kSecAttrAccount as String: "auth_token",
            kSecReturnData as String: true,
            kSecMatchLimit as String: kSecMatchLimitOne
        ]
        var result: AnyObject?
        SecItemCopyMatching(query as CFDictionary, &result)
        guard let data = result as? Data else { return nil }
        return String(data: data, encoding: .utf8)
    }
}
```

**关键点**：
- `kSecAttrAccessible` 设置为 `kSecAttrAccessibleAfterFirstUnlock` 确保设备解锁后 Extension 可读取
- `kSecAttrAccessGroup` 必须设置为 Keychain Group 名称（非 App Group ID）
- 主 App 登录成功后调用 `KeychainService.save(token:)`，Share Extension 中调用 `KeychainService.getToken()` 即可获取

#### 用户无感设计

```
主 App 登录 → 保存 token 到 Keychain (App Group 共享)
        ↓
用户分享 → Share Extension 启动 → 读取 Keychain 获取 token
        ↓
    有 token → 直接调用 API
    无 token → 通知用户打开 App 登录（或通过 openURL 跳转）
```

---

### 4. WKWebView Shell

**最小实现**：一个极薄的 iOS App，主界面仅包含一个全屏 `WKWebView` 加载 React SPA。

#### AppDelegate.swift

```swift
import UIKit

@main
class AppDelegate: UIResponder, UIApplicationDelegate {
    var window: UIWindow?

    func application(_ application: UIApplication,
                     didFinishLaunchingWithOptions launchOptions: [UIApplication.LaunchOptionsKey: Any]?) -> Bool {
        return true
    }
}
```

#### ViewController.swift

```swift
import UIKit
import WebKit

class ViewController: UIViewController, WKUIDelegate, WKNavigationDelegate {

    var webView: WKWebView!

    override func loadView() {
        let config = WKWebViewConfiguration()
        let preferences = WKPreferences()
        preferences.javaScriptEnabled = true
        config.preferences = preferences

        webView = WKWebView(frame: .zero, configuration: config)
        webView.uiDelegate = self
        webView.navigationDelegate = self
        view = webView
    }

    override func viewDidLoad() {
        super.viewDidLoad()
        if let url = URL(string: "https://app.hoard.app") {
            webView.load(URLRequest(url: url))
        }
    }
}
```

#### Info.plist 关键配置

```xml
<key>NSAppTransportSecurity</key>
<dict>
    <key>NSAllowsArbitraryLoads</key>
    <false/>
    <key>NSAllowsArbitraryLoadsInWebContent</key>
    <true/>
</dict>
```

`NSAllowsArbitraryLoadsInWebContent` 仅允许 WebView 内的任意加载，比 `NSAllowsArbitraryLoads` 更安全，App Store 审核更容易通过。

#### URL Scheme / Deep Link 处理

通过 **WKWebView 导航拦截** 实现：

```swift
func webView(_ webView: WKWebView,
             decidePolicyFor navigationAction: WKNavigationAction,
             decisionHandler: @escaping (WKNavigationActionPolicy) -> Void) {

    guard let url = navigationAction.request.url else {
        decisionHandler(.allow)
        return
    }

    if url.scheme == "hoard" {
        handleHoardURL(url)
        decisionHandler(.cancel)
        return
    }

    if url.scheme == "http" || url.scheme == "https" {
        if url.host == "app.hoard.app" {
            decisionHandler(.allow)
        } else {
            UIApplication.shared.open(url, options: [:])
            decisionHandler(.cancel)
        }
        return
    }

    decisionHandler(.allow)
}
```

**JS → Native 通信**：通过 `WKUserContentController` 添加消息处理器：

```swift
config.userContentController.add(self, name: "hoardBridge")

extension ViewController: WKScriptMessageHandler {
    func userContentController(_ userContentController: WKUserContentController,
                                didReceive message: WKScriptMessage) {
        guard let body = message.body as? [String: Any],
              let action = body["action"] as? String else { return }
        switch action {
        case "share":
            if let url = body["url"] as? String {
                presentShareSheet(url: url)
            }
        default:
            break
        }
    }
}
```

**反向通信（Native → JS）**：

```swift
func notifySPA(event: String, data: [String: Any]) {
    let json = (try? JSONSerialization.data(withJSONObject: data))
        .flatMap { String(data: $0, encoding: .utf8) } ?? "{}"
    webView.evaluateJavaScript("window.__hoardBridge?.onEvent('\(event)', \(json))")
}
```

#### 注册自定义 URL Scheme

在 `Info.plist` 中：

```xml
<key>CFBundleURLTypes</key>
<array>
    <dict>
        <key>CFBundleURLSchemes</key>
        <array>
            <string>hoard</string>
        </array>
        <key>CFBundleURLName</key>
        <string>com.hoard.app</string>
    </dict>
</array>
```

在 `SceneDelegate` 中处理：

```swift
func scene(_ scene: UIScene, openURLContexts URLContexts: Set<UIOpenURLContext>) {
    guard let url = URLContexts.first?.url else { return }
    notifySPA(event: "deepLink", data: ["url": url.absoluteString])
}
```

---

### 5. App Store 发布要求

#### 必备条件

| 项目 | 说明 |
|------|------|
| **Apple Developer Program** | 年费 $99，必须注册 |
| **App ID** | 在 Apple Developer Portal 创建，需包含 App Groups / Keychain Sharing 等 Capabilities |
| **证书** | Distribution Certificate（发布用） |
| **Provisioning Profile** | App Store Distribution Profile（含所有 Capabilities） |
| **Xcode Archive** | Product → Archive → Distribute App |

#### 注意事项

1. **Bundle ID 体系**：主 App 和 Extension 的 Bundle ID 需在同一 App ID 体系下，如：
   - 主 App: `com.hoard.app`
   - Share Extension: `com.hoard.app.share-extension`

2. **Capabilities 一致性**：App Groups 和 Keychain Sharing 必须在主 App 和 Extension 的 **Entitlements** 文件中一致配置。

3. **审核注意点**：
   - Share Extension 必须有明确的 UI，不能静默执行（`SLComposeServiceViewController` 天然满足）
   - Extension 不能过度请求权限（只请求必要的 App Group 和 Keychain 访问）
   - WKWebView Shell 必须提供有意义的原生功能——但 Hoard 的收藏夹功能 + 原生分享集成本身已经构成足够的功能性

4. **Xcode Cloud / CI 构建**：

   ```bash
   xcodebuild -workspace Hoard.xcworkspace \
              -scheme Hoard \
              -configuration Release \
              -archivePath ./build/Hoard.xcarchive \
              archive
   ```

---

### 总结：最小实现工作量估算

| 模块 | 文件数 | 预估代码行数 |
|------|--------|-------------|
| Share Extension | 3 个文件（Info.plist + VC + Storyboard） | ~80 行 Swift |
| Keychain Sharing | 1 个工具类 | ~40 行 Swift |
| WKWebView Shell | 2 个文件（AppDelegate + ViewController） | ~100 行 Swift |
| URL Scheme 处理 | 集成在 ViewController 中 | ~30 行 Swift |
| **总计** | **6 个文件** | **~250 行 Swift** |
