# ItdTime ⏰

[English](./README.md) | **简体中文**

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](./LICENSE)
[![Wails](https://img.shields.io/badge/Wails-v3--beta-DF0000?logo=go&logoColor=white)](https://v3.wails.io/)
[![Go](https://img.shields.io/badge/Go-1.26+-00ADD8?logo=go&logoColor=white)](https://go.dev/)
[![Vue](https://img.shields.io/badge/Vue-3.5-4FC08D?logo=vuedotjs&logoColor=white)](https://vuejs.org/)
[![PrimeVue](https://img.shields.io/badge/PrimeVue-4-6366F1?logo=vue.js&logoColor=white)](https://primevue.org/)
[![Platform](https://img.shields.io/badge/Platform-Windows-0078D6?logo=windows&logoColor=white)](#)
[![Download for Windows](https://img.shields.io/badge/Windows_%E4%B8%8B%E8%BD%BD-v1.2.1-2EA043?logo=windows&logoColor=white)](https://wlitd.com/updates/ItdTime-amd64-installer.exe)

> I'm off duty! —— 一个下班提醒 + 待办提醒的桌面小应用
>
> 到点就要准时挂钟离岗 ⏰

`ItdTime` 是一款基于 **Wails v3 + Go + Vue 3** 的轻量桌面应用：设定下班时间后，倒计时归零时会在屏幕右下角弹出置顶提醒卡片，提醒你准时跑路；同时提供待办提醒列表管理零散计划，支持应用内自更新。

## ⬇️ 下载

**[Windows 版下载（v1.2.1，`.exe` 安装器）](https://wlitd.com/updates/ItdTime-amd64-installer.exe)**

需 Windows 10/11 x64。按用户安装到 `%LOCALAPPDATA%\Programs\ItdTime`，免管理员权限；装好后新版本由应用内自更新送达。

## ✨ 功能特性

- **自由倒计时**：LED 数字风格的下班倒计时，实时显示距离下班还有多久
- **到点提醒**：倒计时归零弹出右下角置顶提醒卡片，支持**并行错位堆叠**（多条互不阻塞）
- **提前提醒**：支持设置提前 N 分钟提醒，先收拾收拾再跑路
- **重复规则**：按周一 ~ 周日自由设置重复提醒日
- **提醒样式**：可插拔的样式注册表（当前内置「章鱼哥」），持续扩充中
- **预览提醒**：设置面板内一键预览提醒弹窗效果
- **待办提醒列表**：添加、管理更多计划提醒事项，到期自动弹出
- **系统托盘**：关闭窗口隐藏到托盘，后台持续守候，不错过提醒
- **开机自启**：登录静默启动（后台驻留托盘），可一键开关
- **应用内自更新**：检查 → 下载（进度条）→ 一键替换重启，sha256 完整性校验
- **主题 / 多语言**：深浅色主题切换，中英文双语界面

## 🏗️ 架构

Wails v3 单进程多窗口：提醒弹窗是同进程内的独立透明置顶小窗，由 Go 端编排器直接创建 / 定位 / 关闭。

```
主进程（单 exe）
├── 主窗口（372×600 无边框，前端自绘标题栏）
├── 提醒窗口 ×N（380×64 透明置顶小窗，右下角自底向上错位堆叠，上限 6 条）
├── Go 端 1s 轮询调度 —— 判定逻辑是 domain.NextFires 纯函数（有单测），
│   不受 WebView 后台节流影响，到点即开窗、独立计时 20s 自动关闭
└── 托盘 / 开机自启 / 自更新 —— 均为 Wails v3 原生能力（SystemTray / Autostart / updater）
```

后端按 DDD 轻量分层，**数据主权在 Go 端**（单一 JSON 存储，前端经类型化绑定读写）：

```
internal/
├── api/        # 对前端暴露的绑定服务：settings / todos / system / update
├── domain/     # 领域模型 + 调度纯函数（NextFires）+ 单元测试
├── reminder/   # Coordinator：轮询检测、并行堆叠弹窗、Dismiss/Preview
└── storage/    # 单文件 JSON 持久化（xdg 数据目录）
```

## 🛠️ 技术栈

| 分类 | 技术 |
| --- | --- |
| 桌面框架 | [Wails v3](https://v3.wails.io/)（beta.32）+ [go-task](https://taskfile.dev/) 构建编排 |
| 后端 | Go 1.26 |
| 前端框架 | Vue 3.5 `<script setup>` + Vite 8 + TypeScript 6.0 |
| UI 组件库 | [PrimeVue 4](https://primevue.org/)（自定义主题 preset，深浅双模式） |
| 原子化 CSS | [UnoCSS](https://unocss.dev/)（Lucide / Flag 图标集） |
| 状态管理 | Pinia 4 |
| 国际化 | vue-i18n 11（zh-CN / en） |
| 自动导入 | unplugin-auto-import + unplugin-vue-components |
| 自更新 | Wails v3 updater + endpoint provider（静态 manifest + sha256） |

## 🚀 快速开始

### 环境要求

- [Go](https://go.dev/) 1.26+
- [Node.js](https://nodejs.org/) 20+ 与 [pnpm](https://pnpm.io/)
- Wails v3 CLI：`go install github.com/wailsapp/wails/v3/cmd/wails3@latest`（本项目对齐 `beta.32`）
- Windows 需安装 [WebView2 Runtime](https://developer.microsoft.com/microsoft-edge/webview2/)（Win10/11 一般已内置）

### 开发

```powershell
# 项目根目录（Windows PowerShell；须指定 pnpm）
$env:PACKAGE_MANAGER = 'pnpm'
wails3 task dev          # 热重载，前端端口默认 9245
```

如需单独调试前端（纯浏览器环境，Wails 绑定按 `isWails()` 自动降级为占位态）：

```bash
cd frontend
pnpm install
pnpm dev
# 浏览器直接预览提醒卡片：http://localhost:9245/?remind=1&kind=offwork&title=测试&id=1
```

### 构建 / 打包

```powershell
wails3 task build                              # 开发态构建 → bin/ItdTime.exe
wails3 task package INSTALL_SCOPE=user         # NSIS 用户级安装包 → bin/ItdTime-amd64-installer.exe
```

> `INSTALL_SCOPE=user` 安装到 `%LOCALAPPDATA%\Programs\ItdTime`，免 UAC，且应用内自更新可就地替换。

## 🔄 自更新与发布流程

更新端点固定为 `https://wlitd.com/updates/manifest.json`（[build/updates/manifest.json](./build/updates/manifest.json) 即其发布源）。一键发版（Git Bash）：

```bash
bash scripts/release.sh 1.1.0   # 同步版本号到五处 → 打包 → 复制产物 → 回填 size/digest → 校验
bash scripts/release.sh         # 不升版本，按当前版本重新构建
```

脚本自动化除「上传」外的全部步骤；`manifest.json` 的 `notes`（更新说明，展示在更新弹窗中）需在发版前手动编辑。最终需人工上传两个文件到服务器 `updates/` 目录：`manifest.json` 与 `ItdTime-<version>-windows-amd64.exe`（裸 exe 供应用内更新就地替换，不是安装器）。

手工等价步骤（兜底）：

1. **同步版本号**（四处）：`main.go` 的 `var version`、`build/config.yml`、`build/windows/info.json`、`build/windows/nsis/wails_tools.nsh`
2. 打包：`wails3 task package INSTALL_SCOPE=user`
3. 复制产物：`bin/ItdTime.exe` → `bin/ItdTime-<version>-windows-amd64.exe`
4. 回填 [manifest.json](./build/updates/manifest.json) 的 `size`（产物字节数）与 `digest`（sha256 原始字节的 base64）
5. 将 `manifest.json` 与 exe 上传到服务器 `updates/` 目录

更新落地采用 rename-aside 就地替换（旧 exe 改名 `.old` 后落新文件并重启），启动时自动清理上次残留的 `.old` 文件。

## 🎨 扩充新的提醒样式

1. 在 `frontend/src/features/remind/` 新建样式组件（参考 `SquidwardRemind.vue`）
2. 在 `registry.ts` 的 `resolveRemind` 中注册样式 → 组件映射
3. 样式枚举定义在 Go 端 `internal/domain`（`RemindStyle`），新增样式需同步后端并重新生成绑定

## 📁 项目结构

```
itd-time/
├── main.go                  # 入口：装配 Services、主窗口、托盘
├── Taskfile.yml             # go-task 根编排（dev/build/package 分发到平台 Taskfile）
├── scripts/release.sh       # 一键发版脚本（Git Bash，见「发布流程」）
├── internal/                # Go 后端（api / domain / reminder / storage，见「架构」节）
├── build/
│   ├── config.yml           # Wails v3 构建配置（info 段为版本元数据源）
│   ├── appicon.png          # 应用/托盘图标唯一源（generate:icons 派生 ico/icns）
│   ├── windows/             # exe 版本资源、NSIS 安装器脚本
│   └── updates/             # 自更新 manifest（发布源）
└── frontend/
    ├── DESIGN.md            # 设计系统唯一事实源（色板/字体/间距/组件约定）
    └── src/
        ├── core/            # i18n / PrimeVue 主题 preset / 绑定封装 / 全局 store
        ├── features/        # 按领域切分：countdown / remind / settings / shell / todos / update
        ├── pages/           # MainView 主页面（提醒窗由 App 按 URL 分发）
        ├── shared/          # 运行环境判定、日期工具
        └── styles/          # CSS 变量 + LED 字体
```

## 📄 说明

- 主窗口为固定尺寸（372 × 600）的无边框窗口，通过顶部栏拖拽移动
- 点击窗口关闭按钮不会退出应用，而是隐藏到系统托盘；彻底退出请使用托盘菜单「退出」
- 应用目前面向 Windows 平台（macOS 构建链已随 Wails 保留，未发布）

---

By [Wlitd](mailto:2583651948@qq.com) · 祝大家都能准时下班 🎉
