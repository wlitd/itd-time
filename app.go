package main

import (
	"context"
	_ "embed"
	"os"

	"github.com/wailsapp/wails/v2/pkg/runtime"

	"itd-time/internal/services"
)

//go:embed icon.ico
var iconBytes []byte

// remindWindowTitle 提醒弹窗窗口标题，用于定位窗口句柄
const remindWindowTitle = "ItdTime-remind"

// App 应用主结构体，聚合各服务模块，并作为 Wails 前端绑定的入口
type App struct {
	ctx context.Context

	tray      *services.TrayService
	remind    *services.RemindService
	autoStart *services.AutoStartService
	updater   *services.UpdaterService
}

// NewApp 创建应用实例，初始化各服务模块
func NewApp() *App {
	remindSvc := services.NewRemindService()
	remindSvc.ParseArgs(os.Args[1:])

	return &App{
		tray:      services.NewTrayService(iconBytes),
		remind:    remindSvc,
		autoStart: services.NewAutoStartService(),
		updater:   services.NewUpdaterService(),
	}
}

// startup Wails 启动回调，初始化上下文并分发到各服务
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	a.updater.Startup(ctx)
	a.remind.Start(ctx)

	if a.remind.IsRemind() {
		// 提醒弹窗：移动到屏幕工作区右下角，等待前端渲染完成后再显示
		go services.MoveWindowToBottomRight(remindWindowTitle, 12)
		return
	}

	// 主进程：启动系统托盘（定时器由 remind.Start 内部启动）
	go a.tray.Run(ctx)
}

// OnBeforeClose 窗口关闭前回调
func (a *App) OnBeforeClose(ctx context.Context) (prevent bool) {
	// 提醒弹窗进程直接允许退出
	if a.remind.IsRemind() {
		return false
	}
	if !a.tray.IsQuiting() {
		runtime.WindowHide(ctx)
		return true
	}
	return false
}

// ======================= Wails 前端绑定方法 =======================

// --- 提醒相关 ---

// ShowRemind 启动提醒弹窗子进程
func (a *App) ShowRemind(style string, mode string) {
	a.remind.ShowRemind(style, mode)
}

// SyncRemindConfig 同步提醒配置
func (a *App) SyncRemindConfig(offTime string, advanceMinutes int, repeatDays []int, style string) {
	a.remind.SyncConfig(offTime, advanceMinutes, repeatDays, style)
}

// IsRemindWindow 当前进程是否为提醒弹窗进程
func (a *App) IsRemindWindow() bool {
	return a.remind.IsRemind()
}

// GetRemindStyle 获取提醒弹窗样式
func (a *App) GetRemindStyle() string {
	return a.remind.Style()
}

// GetRemindMode 获取提醒弹窗模式
func (a *App) GetRemindMode() string {
	return a.remind.Mode()
}

// CloseRemind 关闭提醒弹窗
func (a *App) CloseRemind() {
	a.remind.CloseRemind()
}

// --- 自启动相关 ---

// EnableAutoStart 开启开机自启动
func (a *App) EnableAutoStart() error {
	return a.autoStart.EnableAutoStart()
}

// DisableAutoStart 关闭开机自启动
func (a *App) DisableAutoStart() error {
	return a.autoStart.DisableAutoStart()
}

// IsAutoStartEnabled 检查自启动状态
func (a *App) IsAutoStartEnabled() bool {
	return a.autoStart.IsAutoStartEnabled()
}

// --- 更新相关 ---

// GetVersion 获取当前应用版本号
func (a *App) GetVersion() string {
	return services.Version
}

// CheckForUpdate 检查更新
func (a *App) CheckForUpdate() services.UpdateResult {
	return a.updater.CheckForUpdate()
}

// DownloadAndInstall 下载并安装更新
func (a *App) DownloadAndInstall(url string) error {
	return a.updater.DownloadAndInstall(url)
}

// ExitApp 彻底退出应用进程（安装更新前使用，避免托盘残留阻止 exe 覆盖）
func (a *App) ExitApp() {
	os.Exit(0)
}
