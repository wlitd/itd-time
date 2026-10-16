package main

import (
	"embed"
	"log"
	"os"
	"slices"

	"github.com/wailsapp/wails/v3/pkg/application"

	"itdtime/internal/api"
	"itdtime/internal/reminder"
	"itdtime/internal/storage"
)

// 前端产物在构建时嵌入二进制，运行期由内置资源服务器提供。
//
//go:embed all:frontend/dist
var assets embed.FS

// 托盘图标（Windows 通知区用 PNG 字节）。
//
//go:embed build/appicon.png
var trayIconPNG []byte

// autostartArg 与 api.SystemService 约定一致：由系统登录项拉起时携带，静默启动。
const autostartArg = "--autostart"

// version 应用版本号，本变量即运行时唯一真值源（不发版流水线注入）。
// 需为合法 semver（不带 v 前缀），updater 依此与 manifest 版本做新旧比较；
// 发版时须与 build/config.yml、build/windows/info.json、build/windows/nsis/wails_tools.nsh 四处同步。
var version = "1.1.2"

func main() {
	// 基础设施：加载持久化数据。
	store, err := storage.New()
	if err != nil {
		log.Fatalf("init storage: %v", err)
	}
	if err := store.Load(); err != nil {
		log.Fatalf("load data: %v", err)
	}

	// 装配层：构造对前端暴露的 Services。
	settings := api.NewSettingsService(store)
	todos := api.NewTodoService(store)
	system := api.NewSystemService()
	coordinator := reminder.NewCoordinator(store)
	update := api.NewUpdateService(store, version)

	app := application.New(application.Options{
		Name:        "ItdTime",
		Description: "下班倒计时与提醒工具",
		// 应用级图标（窗口/任务栏/Dock）：嵌入 build/appicon.png，替换 Wails 默认图标。
		// 托盘另见 setupTray；Windows exe 文件图标由 build/windows/icon.ico（generate icons 产出）经 syso 嵌入。
		Icon: trayIconPNG,
		Services: []application.Service{
			application.NewService(settings),
			application.NewService(todos),
			application.NewService(system),
			application.NewService(coordinator),
			application.NewService(update),
		},
		Assets: application.AssetOptions{
			Handler: application.AssetFileServerFS(assets),
		},
		Mac: application.MacOptions{
			ApplicationShouldTerminateAfterLastWindowClosed: false,
		},
	})

	// 主窗口：固定尺寸、无边框，由前端自绘标题栏。
	// 开机自启拉起时以隐藏态启动（后台驻留托盘），不抢占前台。
	startHidden := slices.Contains(os.Args, autostartArg)
	app.Window.NewWithOptions(application.WebviewWindowOptions{
		Name:             "main",
		Title:            "ItdTime",
		Width:            372,
		Height:           600,
		Frameless:        true,
		DisableResize:    true,
		BackgroundColour: application.NewRGB(255, 255, 255),
		Hidden:           startHidden,
		URL:              "/",
	})

	// 系统托盘：左键唤出主窗，右键菜单含显示/退出。
	setupTray(app, system)

	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}

// setupTray 创建通知区图标与菜单。
func setupTray(app *application.App, system *api.SystemService) {
	tray := app.SystemTray.New()
	tray.SetIcon(trayIconPNG)
	// 延续 V2 的幽默悬浮文案（原 systray.SetTooltip），而非冷冰冰的应用名。
	tray.SetTooltip("wlitd 正在盯着时钟 ⏰")

	menu := application.NewMenu()
	menu.Add("显示主窗口").OnClick(func(*application.Context) { system.ShowMainWindow() })
	menu.AddSeparator()
	menu.Add("退出").OnClick(func(*application.Context) { system.Quit() })
	tray.SetMenu(menu)

	// 左键单击托盘图标：直接唤出主窗口（关窗即隐藏后靠此恢复到前台）。
	tray.OnClick(func() { system.ShowMainWindow() })

	tray.Run()
}
