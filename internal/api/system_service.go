package api

import (
	"context"
	"os"
	"slices"

	"github.com/wailsapp/wails/v3/pkg/application"
)

// autostartArgName 注册登录项时附带的参数：以此判断「由系统自启拉起」→ 后台静默运行。
const autostartArgName = "--autostart"

// SystemService 暴露系统集成能力：开机自启、主窗显示/隐藏（关窗即隐藏）、退出。
// 它同样通过 ServiceStartup 拿到 App 实例（见 Coordinator 的相同模式）。
type SystemService struct {
	app *application.App
}

// NewSystemService 构造 SystemService。
func NewSystemService() *SystemService {
	return &SystemService{}
}

// ServiceStartup 实现 application.ServiceStartup。
func (s *SystemService) ServiceStartup(_ context.Context, _ application.ServiceOptions) error {
	s.app = application.Get()
	return nil
}

// SetAutostart 开关开机自启；启用时携带 --autostart 参数，令自启拉起时静默运行。
func (s *SystemService) SetAutostart(enabled bool) error {
	if enabled {
		return s.app.Autostart.EnableWithOptions(application.AutostartOptions{
			Arguments: []string{autostartArgName},
		})
	}
	return s.app.Autostart.Disable()
}

// IsAutostartEnabled 查询当前是否已注册开机自启。
func (s *SystemService) IsAutostartEnabled() (bool, error) {
	return s.app.Autostart.IsEnabled()
}

// HideMainWindow 隐藏主窗口（供 frameless 标题栏的「关闭」按钮调用，实现关窗即隐藏）。
func (s *SystemService) HideMainWindow() {
	if w, ok := s.app.Window.GetByName("main"); ok {
		w.Hide()
	}
}

// ShowMainWindow 显示、取消最小化并将主窗口置前（供托盘左键与菜单项复用）。
func (s *SystemService) ShowMainWindow() {
	if w, ok := s.app.Window.GetByName("main"); ok {
		w.UnMinimise()
		w.Show()
		w.Focus()
	}
}

// Quit 退出应用（供托盘菜单与设置页退出项调用）。
func (s *SystemService) Quit() {
	s.app.Quit()
}

// IsLaunchWithAutostart 供前端判断本次进程是否由开机自启拉起（用于静默启动提示）。
func (s *SystemService) IsLaunchWithAutostart() bool {
	return hasAutostartFlag()
}

// hasAutostartFlag 检测命令行是否携带 --autostart（由系统登录项拉起时）。
func hasAutostartFlag() bool {
	return slices.Contains(os.Args, autostartArgName)
}
