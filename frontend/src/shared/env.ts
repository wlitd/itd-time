// 运行环境判定：Wails 桌面运行时 vs 纯浏览器（vite 独立预览）。
// 纯浏览器下所有 bindings 调用不可用，feature 层据此降级为占位/内存态。

/**
 * 是否运行在 Wails 桌面环境。
 * v3 不同平台的桥接注入点各不相同（见 @wailsio/runtime 的 system.js）：
 *   - Windows WebView2：window.chrome.webview.postMessage
 *   - macOS/iOS WKWebView：window.webkit.messageHandlers.external.postMessage
 *   - Android：window.wails.invoke
 * 需与运行时 transport 选择逻辑保持一致；`window._wails.environment` 是异步填充的，
 * 挂载时不可靠，不能用作同步判定。纯浏览器下三者皆无 → false。
 */
export function isWails(): boolean {
  const w = window as any
  return !!(
    w.chrome?.webview?.postMessage
    || w.webkit?.messageHandlers?.external?.postMessage
    || w.wails?.invoke
  )
}

/** 当前窗口是否为提醒窗：由 Coordinator 打开时注入 ?remind=1。 */
export function isRemindWindow(): boolean {
  return new URLSearchParams(window.location.search).get('remind') === '1'
}

/** 读取提醒窗上下文参数（仅提醒窗有效）。 */
export function readRemindParams(): { kind: string, title: string, todoId: string, id: number } {
  const q = new URLSearchParams(window.location.search)
  return {
    kind: q.get('kind') ?? '',
    title: q.get('title') ?? '',
    todoId: q.get('todoId') ?? '',
    // id 由 Coordinator 开窗时注入，供 Dismiss(id) 精确定位并关闭对应提醒窗。
    id: Number(q.get('id') ?? '0'),
  }
}
