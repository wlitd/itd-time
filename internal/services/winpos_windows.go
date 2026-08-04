//go:build windows

package services

import (
	"syscall"
	"time"
	"unsafe"
)

var (
	user32                    = syscall.NewLazyDLL("user32.dll")
	procFindWindowW           = user32.NewProc("FindWindowW")
	procGetWindowRect         = user32.NewProc("GetWindowRect")
	procSetWindowPos          = user32.NewProc("SetWindowPos")
	procSystemParametersInfoW = user32.NewProc("SystemParametersInfoW")
)

// winRect Win32 RECT 结构体
type winRect struct {
	Left, Top, Right, Bottom int32
}

const (
	spiGetWorkArea = 0x0030
	swpNoSize      = 0x0001
	swpNoActivate  = 0x0010
)

// hwndTopmost 置顶窗口层级（HWND_TOPMOST = -1）
var hwndTopmost = ^uintptr(0)

// MoveWindowToBottomRight 将指定标题的窗口移动到主屏幕工作区（不含任务栏）右下角
// margin 为距离屏幕边缘的间距（物理像素）
func MoveWindowToBottomRight(title string, margin int32) {
	ptr, err := syscall.UTF16PtrFromString(title)
	if err != nil {
		return
	}

	// 窗口创建与标题设置存在时序差，重试查找窗口句柄
	var hwnd uintptr
	for range 40 {
		hwnd, _, _ = procFindWindowW.Call(0, uintptr(unsafe.Pointer(ptr)))
		if hwnd != 0 {
			break
		}
		time.Sleep(50 * time.Millisecond)
	}
	if hwnd == 0 {
		return
	}

	// 窗口实际尺寸与工作区均为物理像素，计算结果不受 DPI 缩放影响
	var wr, wa winRect
	procGetWindowRect.Call(hwnd, uintptr(unsafe.Pointer(&wr)))
	procSystemParametersInfoW.Call(spiGetWorkArea, 0, uintptr(unsafe.Pointer(&wa)), 0)

	w := wr.Right - wr.Left
	h := wr.Bottom - wr.Top
	x := wa.Right - w - margin
	y := wa.Bottom - h - margin

	procSetWindowPos.Call(hwnd, hwndTopmost, uintptr(x), uintptr(y), 0, 0, swpNoSize|swpNoActivate)
}
