//go:build !windows && !darwin

package services

// MoveWindowToBottomRight 将窗口移动到屏幕右下角（非 Windows 平台暂不实现）
func MoveWindowToBottomRight(title string, margin int32) {}
