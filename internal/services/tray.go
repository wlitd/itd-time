package services

import (
	"context"

	"fyne.io/systray"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// TrayService 系统托盘服务，管理托盘图标与菜单交互
type TrayService struct {
	ctx       context.Context
	isQuiting bool
	iconBytes []byte
}

// NewTrayService 创建托盘服务
func NewTrayService(iconBytes []byte) *TrayService {
	return &TrayService{iconBytes: iconBytes}
}

// IsQuiting 是否正在退出（由退出菜单项触发）
func (s *TrayService) IsQuiting() bool {
	return s.isQuiting
}

// Run 启动系统托盘（阻塞，应在 goroutine 中调用）
func (s *TrayService) Run(ctx context.Context) {
	s.ctx = ctx
	systray.Run(s.onReady, s.onExit)
}

func (s *TrayService) onReady() {
	systray.SetIcon(s.iconBytes)
	systray.SetTitle("ItdTime")
	systray.SetTooltip("wlitd 正在盯着时钟 ⏰")

	mShow := systray.AddMenuItem("显示窗口", "显示主窗口")
	systray.AddSeparator()
	mQuit := systray.AddMenuItem("退出应用", "彻底退出程序")

	go func() {
		for {
			select {
			case <-mShow.ClickedCh:
				runtime.WindowShow(s.ctx)
			case <-mQuit.ClickedCh:
				s.isQuiting = true
				systray.Quit()
				runtime.Quit(s.ctx)
			}
		}
	}()
}

func (s *TrayService) onExit() {}
