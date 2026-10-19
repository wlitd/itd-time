# ItdTime ⏰

**English** | [简体中文](./README-zh.md)

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](./LICENSE)
[![Wails](https://img.shields.io/badge/Wails-v3--beta-DF0000?logo=go&logoColor=white)](https://v3.wails.io/)
[![Go](https://img.shields.io/badge/Go-1.26+-00ADD8?logo=go&logoColor=white)](https://go.dev/)
[![Vue](https://img.shields.io/badge/Vue-3.5-4FC08D?logo=vuedotjs&logoColor=white)](https://vuejs.org/)
[![PrimeVue](https://img.shields.io/badge/PrimeVue-4-6366F1?logo=vue.js&logoColor=white)](https://primevue.org/)
[![Platform](https://img.shields.io/badge/Platform-Windows-0078D6?logo=windows&logoColor=white)](#)
[![Download for Windows](https://img.shields.io/badge/Download_for_Windows-v1.2.0-2EA043?logo=windows&logoColor=white)](https://wlitd.com/updates/ItdTime-amd64-installer.exe)

> I'm off duty! — a tiny desktop app for off-work countdown & to-do reminders
>
> Clock out on time, every time ⏰

`ItdTime` is a lightweight desktop app built with **Wails v3 + Go + Vue 3**: set your off-work time and a countdown tells you how long until freedom; when it hits zero, a topmost reminder card pops up at the bottom-right of the screen to send you home on time. It also keeps a to-do reminder list for scattered plans, and supports in-app self-update.

## ⬇️ Download

**[Download for Windows (v1.2.0, `.exe` installer)](https://wlitd.com/updates/ItdTime-amd64-installer.exe)**

Requires Windows 10/11 x64. Installs per-user into `%LOCALAPPDATA%\Programs\ItdTime` — no admin rights needed. After that, newer versions are delivered by the built-in in-app self-update.

## ✨ Features

- **Countdown**: LED-style off-work countdown showing the time remaining
- **On-time reminder**: a topmost card pops up at the bottom-right when it hits zero; multiple cards **stack in parallel with offsets** (never block each other)
- **Early reminder**: get notified N minutes ahead so you can wrap up first
- **Repeat days**: freely pick which weekdays (Mon–Sun) to remind
- **Reminder styles**: a pluggable style registry (currently "Squidward"), more coming
- **Preview**: one-click preview of the reminder popup right in the settings panel
- **To-do reminders**: add and manage extra planned reminders that pop up when due
- **System tray**: closing the window hides it to the tray, watching on in the background
- **Autostart**: silently launches at login (resides in the tray), toggleable in one click
- **In-app self-update**: check → download (with progress) → replace & restart in one click, verified by sha256
- **Theme / i18n**: light & dark themes, Chinese / English UI

## 🏗️ Architecture

Single process, multiple windows (Wails v3): reminder popups are frameless transparent always-on-top windows in the same process, created / positioned / closed directly by a Go-side orchestrator.

```
Main process (single exe)
├── Main window (372×600 frameless, custom titlebar drawn by the frontend)
├── Reminder windows ×N (380×64 transparent topmost cards, stacked bottom-up
│   at the bottom-right with offsets, capped at 6)
├── Go-side 1s polling scheduler — the decision logic is a pure function
│   (domain.NextFires, unit-tested), immune to WebView background throttling;
│   each card closes automatically after 20s on its own timer
└── Tray / autostart / self-update — all native Wails v3 capabilities
    (SystemTray / Autostart / updater)
```

The backend uses a lightweight DDD layering, with **data sovereignty on the Go side** (a single JSON store; the frontend reads/writes through typed bindings):

```
internal/
├── api/        # Services bound to the frontend: settings / todos / system / update
├── domain/     # Domain models + pure scheduling function (NextFires) + unit tests
├── reminder/   # Coordinator: polling, parallel stacked popups, Dismiss/Preview
└── storage/    # Single-file JSON persistence (xdg data directory)
```

## 🛠️ Tech Stack

| Area | Technology |
| --- | --- |
| Desktop framework | [Wails v3](https://v3.wails.io/) (beta.32) + [go-task](https://taskfile.dev/) orchestration |
| Backend | Go 1.26 |
| Frontend | Vue 3.5 `<script setup>` + Vite 8 + TypeScript 6.0 |
| UI library | [PrimeVue 4](https://primevue.org/) (custom theme preset, light & dark) |
| Atomic CSS | [UnoCSS](https://unocss.dev/) (Lucide / Flag icon sets) |
| State | Pinia 4 |
| i18n | vue-i18n 11 (zh-CN / en) |
| Auto-import | unplugin-auto-import + unplugin-vue-components |
| Self-update | Wails v3 updater + endpoint provider (static manifest + sha256) |

## 🚀 Getting Started

### Prerequisites

- [Go](https://go.dev/) 1.26+
- [Node.js](https://nodejs.org/) 20+ with [pnpm](https://pnpm.io/)
- Wails v3 CLI: `go install github.com/wailsapp/wails/v3/cmd/wails3@latest` (this project targets `beta.32`)
- [WebView2 Runtime](https://developer.microsoft.com/microsoft-edge/webview2/) on Windows (usually preinstalled on Win10/11)

### Development

```powershell
# At the repository root (Windows PowerShell; pnpm must be specified)
$env:PACKAGE_MANAGER = 'pnpm'
wails3 task dev          # hot reload, frontend dev port defaults to 9245
```

To debug the frontend standalone (pure browser environment; Wails binding calls degrade to placeholders via `isWails()`):

```bash
cd frontend
pnpm install
pnpm dev
# Preview a reminder card directly in the browser:
# http://localhost:9245/?remind=1&kind=offwork&title=Test&id=1
```

### Build / Package

```powershell
wails3 task build                              # dev-mode build → bin/ItdTime.exe
wails3 task package INSTALL_SCOPE=user         # NSIS per-user installer → bin/ItdTime-amd64-installer.exe
```

> `INSTALL_SCOPE=user` installs into `%LOCALAPPDATA%\Programs\ItdTime` without UAC, which also lets in-app self-update replace the binary in place.

## 🔄 Self-Update & Release Flow

The update endpoint is fixed at `https://wlitd.com/updates/manifest.json` ([build/updates/manifest.json](./build/updates/manifest.json) is its publishing source). One-command release (Git Bash):

```bash
bash scripts/release.sh 1.1.0   # sync version across 5 places → package → copy artifact → backfill size/digest → verify
bash scripts/release.sh         # rebuild without bumping the version
```

The script automates everything except the upload itself; edit `notes` in `manifest.json` (shown in the update dialog) manually before releasing. Two files must be uploaded to the server's `updates/` directory: `manifest.json` and `ItdTime-<version>-windows-amd64.exe` (the raw exe is what in-app updates swap in — not the installer).

Manual equivalent (fallback):

1. **Sync the version** (four places): `var version` in `main.go`, `build/config.yml`, `build/windows/info.json`, `build/windows/nsis/wails_tools.nsh`
2. Package: `wails3 task package INSTALL_SCOPE=user`
3. Copy the artifact: `bin/ItdTime.exe` → `bin/ItdTime-<version>-windows-amd64.exe`
4. Backfill `size` (exact byte count of the artifact) and `digest` (base64 of the raw sha256 bytes) in [manifest.json](./build/updates/manifest.json)
5. Upload both `manifest.json` and the exe to the server's `updates/` directory

Updates apply via rename-aside in place (the running exe is renamed to `.old`, the new file lands, then the app restarts); stale `.old` files are cleaned up automatically on next startup.

## 🎨 Adding a New Reminder Style

1. Create a style component under `frontend/src/features/remind/` (see `SquidwardRemind.vue`)
2. Register the style → component mapping in `resolveRemind` inside `registry.ts`
3. The style enum lives in the Go backend (`internal/domain`, `RemindStyle`) — new styles require a backend update and regenerated bindings

## 📁 Project Structure

```
itd-time/
├── main.go                  # Entry: assembles Services, main window, tray
├── Taskfile.yml             # go-task root orchestration (dev/build/package dispatch per platform)
├── scripts/release.sh       # one-command release script (Git Bash, see Release Flow)
├── internal/                # Go backend (api / domain / reminder / storage, see Architecture)
├── build/
│   ├── config.yml           # Wails v3 build config (info section = version metadata source)
│   ├── appicon.png          # Single source for app/tray icons (generates ico/icns)
│   ├── windows/             # exe version resources, NSIS installer scripts
│   └── updates/             # self-update manifest (publishing source)
└── frontend/
    ├── DESIGN.md            # Design system — single source of truth
    └── src/
        ├── core/            # i18n / PrimeVue theme preset / bindings wrapper / global store
        ├── features/        # By domain: countdown / remind / settings / shell / todos / update
        ├── pages/           # MainView (reminder windows are dispatched by App via URL)
        ├── shared/          # Runtime environment detection, date utils
        └── styles/          # CSS variables + LED font
```

## 📄 Notes

- The main window is a fixed-size (372 × 600) frameless window, moved by dragging its top bar
- Clicking the close button hides the app to the system tray instead of quitting; use "Exit" in the tray menu to quit fully
- The app currently targets Windows (the macOS build chain is kept via Wails but not released)

---

By [Wlitd](mailto:2583651948@qq.com) · May everyone clock out on time 🎉
