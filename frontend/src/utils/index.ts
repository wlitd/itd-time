export function getEnv(key: string): string {
  const env = import.meta.env

  return env[key]
}

/**
 * 判断当前是否运行在 Wails 桌面环境中（浏览器调试时为 false）
 */
export function isWails(): boolean {
  return typeof window !== 'undefined' && 'go' in window
}
