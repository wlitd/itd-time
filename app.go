package main

import (
	"context"
	_ "embed"
	"os"
	"strings"
	"sync"

	"fyne.io/systray"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

//go:embed icon.ico
var iconBytes []byte

// remindWindowTitle 提醒弹窗窗口标题，用于定位窗口句柄
const remindWindowTitle = "itd-time-remind"

// App struct
type App struct {
	ctx       context.Context
	isQuiting bool

	// 提醒弹窗进程相关（通过命令行参数 --remind 启动）
	isRemind    bool
	remindStyle string
	remindMode  string

	// 提醒配置（由前端同步，主进程的定时器据此判断提醒时机）
	remindMu  sync.Mutex
	remindCfg remindConfig
}

// NewApp creates a new App application struct
func NewApp() *App {
	a := &App{}

	// 解析命令行参数，判断是否为提醒弹窗进程
	for _, arg := range os.Args[1:] {
		switch {
		case arg == "--remind":
			a.isRemind = true
		case strings.HasPrefix(arg, "--style="):
			a.remindStyle = strings.TrimPrefix(arg, "--style=")
		case strings.HasPrefix(arg, "--mode="):
			a.remindMode = strings.TrimPrefix(arg, "--mode=")
		}
	}

	return a
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	if a.isRemind {
		// 提醒弹窗：移动到屏幕工作区右下角，等待前端渲染完成后再显示
		go moveWindowToBottomRight(remindWindowTitle, 12)
		return
	}

	go systray.Run(a.onTrayReady, a.onTrayExit)
	go a.runRemindTicker()
}

// onTrayReady is called when the systray is ready
func (a *App) onTrayReady() {
	systray.SetIcon(iconBytes)
	systray.SetTitle("itdTime")
	systray.SetTooltip("wlitd 正在盯着时钟 ⏰")

	mShow := systray.AddMenuItem("显示窗口", "显示主窗口")
	systray.AddSeparator()
	mQuit := systray.AddMenuItem("退出应用", "彻底退出程序")

	go func() {
		for {
			select {
			case <-mShow.ClickedCh:
				runtime.WindowShow(a.ctx)
			case <-mQuit.ClickedCh:
				a.isQuiting = true
				systray.Quit()
				runtime.Quit(a.ctx)
			}
		}
	}()
}

func (a *App) onTrayExit() {}

func (a *App) OnBeforeClose(ctx context.Context) (prevent bool) {
	// 提醒弹窗进程直接允许退出
	if a.isRemind {
		return false
	}
	if !a.isQuiting {
		runtime.WindowHide(ctx)
		return true
	}
	return false
}
