# itd-time ⏰

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](./LICENSE)
[![Wails](https://img.shields.io/badge/Wails-v2-DF0000?logo=wails&logoColor=white)](https://wails.io/)
[![Go](https://img.shields.io/badge/Go-1.23+-00ADD8?logo=go&logoColor=white)](https://go.dev/)
[![Vue](https://img.shields.io/badge/Vue-3-4FC08D?logo=vuedotjs&logoColor=white)](https://vuejs.org/)
[![Platform](https://img.shields.io/badge/Platform-Windows-0078D6?logo=windows&logoColor=white)](#)

> I'm off duty! —— 一个下班提醒 + 待办提醒的桌面小应用
>
> 到点就要准时挂钟离岗 ⏰

`itd-time` 是一款基于 **Wails + Go + Vue 3** 的轻量桌面应用，核心功能是下班倒计时与到点提醒：设定下班时间后，应用会在屏幕右下角弹出置顶提醒弹窗，提醒你准时跑路；同时还提供一个待办提醒列表，管理零散的计划事项。

## ✨ 功能特性

- **自由倒计时**：LED 数字风格的下班倒计时，实时显示距离下班还有多久
- **到点提醒**：倒计时归零时，屏幕右下角弹出置顶提醒弹窗
- **提前提醒**：支持设置提前 N 分钟提醒，先收拾收拾再跑路
- **重复规则**：按周一 ~ 周日自由设置重复提醒日
- **提醒样式**：提醒弹窗支持多种样式，自由选择、持续扩充中
- **预览提醒**：设置面板内一键预览提醒弹窗效果
- **待办提醒列表**：添加、管理更多计划提醒事项
- **系统托盘**：关闭窗口自动隐藏到托盘，后台持续守候，不错过提醒
- **主题 / 多语言**：深浅色主题切换，中英文双语界面

## 🖥️ 提醒弹窗

提醒弹窗的基础交互：

- 显示在屏幕**工作区右下角**（自动避开任务栏，适配 DPI 缩放）
- 无边框、置顶、半透明圆角卡片，支持拖拽
- 弹窗内容由所选提醒样式组件渲染，各样式可自由发挥
- 30 秒无操作自动关闭

由于 Wails v2 仅支持单窗口，提醒弹窗采用**独立进程**方案实现：

```
主进程（主窗口 + 托盘 + Go 秒级定时器）
   │  到达提醒时间点
   ▼
exec 自身 --remind --style=xxx --mode=xxx
   │
   ▼
提醒进程（透明置顶小窗，Win32 API 定位右下角）
```

- 提醒时机判断放在 **Go 端定时器**（`remind.go`），不受 WebView 后台节流影响
- 前端配置（下班时间 / 提前分钟数 / 重复日 / 样式）变更时实时同步到 Go 端
- 提醒模式：`advance` 提前提醒 / `offwork` 下班提醒 / `preview` 预览（预览按下班提醒展示）

## 🛠️ 技术栈

| 分类 | 技术 |
| --- | --- |
| 桌面框架 | [Wails v2](https://wails.io/)（Go 1.23） |
| 前端框架 | Vue 3 + Vite 8 + TypeScript |
| UI 组件库 | [Naive UI](https://www.naiveui.com/) |
| 原子化 CSS | [UnoCSS](https://unocss.dev/)（Lucide 图标） |
| 状态管理 | Pinia + pinia-plugin-persistedstate |
| 国际化 | vue-i18n（SFC `<i18n>` 块，中英双语） |
| 系统托盘 | [fyne.io/systray](https://github.com/fyne-io/systray) |
| 工具库 | VueUse / dayjs |

## 🚀 快速开始

### 环境要求

- [Go](https://go.dev/) 1.23+
- [Node.js](https://nodejs.org/) 18+ 与 [pnpm](https://pnpm.io/)
- [Wails CLI](https://wails.io/docs/gettingstarted/installation)：`go install github.com/wailsapp/wails/v2/cmd/wails@latest`
- Windows 需安装 [WebView2 Runtime](https://developer.microsoft.com/microsoft-edge/webview2/)（Win10/11 一般已内置）

### 开发

```bash
# 在项目根目录启动开发模式（热重载）
wails dev
```

如需单独调试前端（浏览器环境，Wails 绑定自动降级为空操作）：

```bash
cd frontend
pnpm install
pnpm dev
```

> 浏览器调试提醒弹窗样式：访问 `/remind?style=squidward&mode=offwork`（`mode` 可选 `advance` / `offwork` / `preview`）

### 构建

```bash
# 构建生产版本，产物位于 build/bin/
wails build
```

## 📁 项目结构

```
itd-time/
├── main.go                 # 应用入口，主窗口 / 提醒弹窗两种启动配置
├── app.go                  # App 结构体：启动参数解析、托盘、生命周期
├── remind.go               # 提醒核心：Go 端定时器、提醒进程拉起、前端绑定方法
├── winpos_windows.go       # Win32 API 将提醒弹窗定位到工作区右下角
├── wails.json              # Wails 项目配置
└── frontend/
    ├── wailsjs/            # Wails 自动生成的 Go 绑定与运行时
    └── src/
        ├── components/
        │   ├── reminds/    # 提醒弹窗样式组件（扩充新样式放这里）
        │   ├── ui/         # 浮动标签风格的基础组件（输入框、选择器等）
        │   ├── CountdownCard.vue    # 倒计时卡片 + 下班提醒设置弹窗
        │   ├── RemindList.vue       # 待办提醒列表
        │   └── RemindStyleSelect.vue# 提醒样式选择器
        ├── composables/
        │   └── useOffWorkRemind.ts  # 提醒配置同步 + 手动弹出提醒
        ├── pages/
        │   ├── index.vue   # 主页面
        │   └── remind.vue  # 提醒弹窗页面（样式分发、自动关闭）
        ├── router/         # 路由与守卫（提醒进程重定向到 /remind）
        └── stores/         # Pinia 状态（下班设置、待办、主题、语言）
```

## 🎨 扩充新的提醒样式

1. 在 `frontend/src/components/reminds/` 新建样式组件（接收 `mode` prop，发出 `close` 事件），可参考 `SquidwardRemind.vue`
2. 在 `frontend/src/stores/offWork.ts` 的 `RemindStyle` 类型中追加样式标识
3. 在 `frontend/src/pages/remind.vue` 的 `remindComponents` 映射中注册组件
4. 在 `frontend/src/components/RemindStyleSelect.vue` 的下拉选项中添加该样式

## 📄 说明

- 主窗口为固定尺寸（348 × 512）的无边框窗口，通过顶部栏拖拽移动
- 点击窗口关闭按钮不会退出应用，而是隐藏到系统托盘；如需彻底退出，请使用托盘菜单的「退出应用」
- 应用目前面向 Windows 平台（系统托盘、提醒弹窗右下角定位依赖 Win32 API）

---

By [Wlitd](mailto:2583651948@qq.com) · 祝大家都能准时下班 🎉
