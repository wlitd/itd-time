package main

import (
	"embed"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/windows"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	// Create an instance of the app structure
	app := NewApp()

	appOptions := &options.App{
		Title:  "itd-time",
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

	// 提醒弹窗进程：右下角置顶小窗，隐藏启动，由前端渲染完成后再显示
	if app.isRemind {
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

	// Create application with options
	err := wails.Run(appOptions)

	if err != nil {
		println("Error:", err.Error())
	}
}
