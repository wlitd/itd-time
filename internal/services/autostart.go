package services

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"syscall"
)

const AppName = "ItdTime"

// AutoStartService 自启动服务结构体
type AutoStartService struct{}

// NewAutoStartService 构造函数
func NewAutoStartService() *AutoStartService {
	return &AutoStartService{}
}

// Enable 开启
func (a *AutoStartService) EnableAutoStart() error {
	switch runtime.GOOS {
	case "windows":
		return a.enableAutoStartWindows()
	case "darwin":
		return a.enableAutoStartMac()
	case "linux":
		return a.enableAutoStartLinux()
	default:
		return fmt.Errorf("不支持当前系统")
	}
}

// DisableAutoStart 关闭自启动
func (a *AutoStartService) DisableAutoStart() error {
	switch runtime.GOOS {
	case "windows":
		return a.disableAutoStartWindows()
	case "darwin":
		return a.disableAutoStartMac()
	case "linux":
		return a.disableAutoStartLinux()
	default:
		return fmt.Errorf("不支持当前系统")
	}
}

// IsAutoStartEnabled 检查自启动状态
func (a *AutoStartService) IsAutoStartEnabled() bool {
	switch runtime.GOOS {
	case "windows":
		return a.isAutoStartEnabledWindows()
	case "darwin":
		_, err := os.Stat(a.getMacPlistPath())
		return err == nil
	case "linux":
		_, err := os.Stat(a.getLinuxDesktopPath())
		return err == nil
	}
	return false
}

// ================= Windows 实现 =================
// 使用注册表方式

func (a *AutoStartService) enableAutoStartWindows() error {
	exePath, _ := os.Executable()
	// 注册表值追加 --autostart 参数，启动时检测该标志以隐藏窗口（后台启动到托盘）
	key := `HKCU\Software\Microsoft\Windows\CurrentVersion\Run`
	cmd := exec.Command("reg", "add", key, "/v", AppName, "/t", "REG_SZ", "/d", exePath+" --autostart", "/f")
	cmd.SysProcAttr = &syscall.SysProcAttr{
		HideWindow: true,
	}
	return cmd.Run()
}

func (a *AutoStartService) disableAutoStartWindows() error {
	key := `HKCU\Software\Microsoft\Windows\CurrentVersion\Run`
	cmd := exec.Command("reg", "delete", key, "/v", AppName, "/f")
	cmd.SysProcAttr = &syscall.SysProcAttr{
		HideWindow: true,
	}
	return cmd.Run()
}

func (a *AutoStartService) isAutoStartEnabledWindows() bool {
	key := `HKCU\Software\Microsoft\Windows\CurrentVersion\Run`
	// 使用 reg query 查询
	cmd := exec.Command("reg", "query", key, "/v", AppName)
	cmd.SysProcAttr = &syscall.SysProcAttr{
		HideWindow: true,
	}
	err := cmd.Run()
	return err == nil
}

// ================= MacOS 实现 =================
// 使用 LaunchAgent (plist 文件)

func (a *AutoStartService) getMacPlistPath() string {
	return filepath.Join(os.Getenv("HOME"), "Library", "LaunchAgents", AppName+".plist")
}

func (a *AutoStartService) enableAutoStartMac() error {
	exePath, _ := os.Executable()
	// Mac 的 plist 文件内容
	content := fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>Label</key>
    <string>%s</string>
    <key>ProgramArguments</key>
    <array>
        <string>%s</string>
        <string>--autostart</string>
    </array>
    <key>RunAtLoad</key>
    <true/>
</dict>
</plist>`, AppName, exePath)

	plistPath := a.getMacPlistPath()
	return os.WriteFile(plistPath, []byte(content), 0644)
}

func (a *AutoStartService) disableAutoStartMac() error {
	plistPath := a.getMacPlistPath()
	if _, err := os.Stat(plistPath); err == nil {
		return os.Remove(plistPath)
	}
	return nil
}

// ================= Linux 实现 =================
// 使用 .desktop 文件 (标准方式)

func (a *AutoStartService) getLinuxDesktopPath() string {
	return filepath.Join(os.Getenv("HOME"), ".config", "autostart", AppName+".desktop")
}

func (a *AutoStartService) enableAutoStartLinux() error {
	exePath, _ := os.Executable()
	content := fmt.Sprintf(`[Desktop Entry]
Type=Application
Name=%s
Exec=%s --autostart
Terminal=false
`, AppName, exePath)

	desktopPath := a.getLinuxDesktopPath()
	// 确保目录存在
	os.MkdirAll(filepath.Dir(desktopPath), 0755)
	return os.WriteFile(desktopPath, []byte(content), 0755)
}

func (a *AutoStartService) disableAutoStartLinux() error {
	desktopPath := a.getLinuxDesktopPath()
	if _, err := os.Stat(desktopPath); err == nil {
		return os.Remove(desktopPath)
	}
	return nil
}
