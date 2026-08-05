package main

import (
	"embed"
	"os"
	"slices"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/windows"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	app := NewApp()

	// 主窗口默认配置
	appOptions := &options.App{
		Title:  "ItdTime",
		Width:  348,
		Height: 512,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},

		Frameless:        true,
		DisableResize:    true,
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup:        app.startup,
		OnBeforeClose:    app.OnBeforeClose,
		Bind: []any{
			app,
		},
	}

	// 提醒弹窗进程：覆盖为右下角置顶小窗配置
	if app.remind.IsRemind() {
		appOptions.Title = remindWindowTitle
		appOptions.Width = 400
		appOptions.Height = 84
		appOptions.AlwaysOnTop = true
		appOptions.StartHidden = true
		// 窗口背景全透明，由前端半透明卡片实现圆角与透明度效果
		appOptions.BackgroundColour = &options.RGBA{R: 0, G: 0, B: 0, A: 0}
		appOptions.Windows = &windows.Options{
			WebviewIsTransparent: true,
			WindowIsTranslucent:  true,
			BackdropType:         windows.None,
		}
	}

	// 自启动（后台启动到托盘）：窗口隐藏，仅显示系统托盘图标
	if hasAutostartFlag() {
		appOptions.StartHidden = true
	}

	err := wails.Run(appOptions)
	if err != nil {
		println("Error:", err.Error())
	}
}

// hasAutostartFlag 检测命令行是否包含 --autostart 参数
func hasAutostartFlag() bool {
	return slices.Contains(os.Args[1:], "--autostart")
}
