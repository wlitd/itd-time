# ItdTime 前端

ItdTime 桌面应用的前端部分，基于 **Vue 3 + Vite 8 + TypeScript** 构建，由 Wails 加载运行。项目整体介绍、桌面端架构与构建发布请参阅[根目录 README](../README.md)。

## 📦 常用命令

日常开发直接在项目根目录运行 `wails dev`（会自动拉起本目录的 Vite 开发服务并注入 Go 绑定）。以下命令用于单独调试前端：

```bash
pnpm install        # 安装依赖（pnpm workspace + catalog 统一管理版本）
pnpm dev            # 启动 Vite 开发服务（纯浏览器环境）
pnpm build          # 构建产物到 dist/（wails build 时自动调用）

npx vue-tsc --build # TypeScript 类型检查
npx eslint src      # 代码规范检查（@antfu/eslint-config）
```

> **浏览器环境说明**：脱离 Wails 运行时，`window.go` 不存在，所有 Go 绑定调用都会经 `isWails()` 判断自动降级为空操作，页面可正常预览。

## 🧱 技术组成

- **UI 组件库**：Naive UI（`unplugin-vue-components` 按需自动注册）
- **原子化 CSS**：UnoCSS + Lucide 图标（`i-lucide:xxx` 类名即图标）
- **状态管理**：Pinia + `pinia-plugin-persistedstate`（持久化到 localStorage）
- **国际化**：vue-i18n，文案写在各 SFC 的 `<i18n lang="json">` 块中（中英双语）
- **路由**：vue-router（history 模式）
- **工具**：VueUse / dayjs

## 📁 目录结构

```
src/
├── assets/            # 字体（LED 数字字体）与图片资源
├── components/
│   ├── layout/        # 布局组件（顶部栏：拖拽区、主题/语言切换、窗口控制）
│   ├── reminds/       # 提醒弹窗样式组件（每种样式一个组件）
│   ├── ui/            # 浮动标签风格基础组件（float-input / float-select 等）
│   └── *.vue          # 业务组件（倒计时卡片、待办列表、重复设置等）
├── composables/       # 组合式函数（useOffWorkRemind：提醒配置同步与弹出）
├── layouts/           # 页面布局（default：顶部栏 + 内容区）
├── pages/
│   ├── index.vue      # 主页面
│   └── remind.vue     # 提醒弹窗页面（样式分发、自动关闭）
├── plugins/           # 应用插件注册（pinia / router / i18n / naive 脱离上下文 API）
├── router/            # 路由与守卫（remindGuard：提醒进程重定向到 /remind）
├── stores/            # Pinia stores（offWork / todo / theme / locale）
├── styles/            # 全局样式与 CSS 变量
└── utils/             # 工具函数（getEnv / isWails）

wailsjs/               # Wails 生成的 Go 绑定与运行时（需提交，供类型检查）
```

## 📐 开发约定

- **自动导入**：`vue` / `vueuse` / `pinia` / `vue-router` 的 API，以及 `composables`、`stores`、`utils` 等目录下的导出均由 `unplugin-auto-import` 自动导入，无需手写 import（声明见 `src/auto-imports.d.ts`）
- **组件自动注册**：`src/components` 下的组件与 Naive UI 组件在模板中直接使用（声明见 `src/components.d.ts`）
- **注释与文案**：代码注释以中文为主；用户可见文案一律写入 SFC 的 `<i18n>` 块并提供中英双语
- **间距规范**：提醒弹窗内横向间距用 `gap-2 / px-2`（8px），纵向用 `gap-3 / py-3`（12px）
- **长文本截断**：使用 `NEllipsis`（带 tooltip），flex 布局下配合 `min-w-0`
- **窗口尺寸**：主窗口固定 348 × 512、提醒弹窗 400 × 84，均为无边框窗口；开发 UI 时注意小尺寸下的显示与滚动行为
- **弹窗内禁用 tooltip**：提醒弹窗等小窗口内不使用 NEllipsis / NTooltip（弹层会打乱布局），文本溢出用纯 CSS `truncate` 兜底，并优先保证标题完整显示
- **窗口拖拽**：通过 CSS 变量 `--wails-draggable: drag / no-drag` 控制可拖拽区域，交互元素需标记 `no-drag`

## 🔔 提醒弹窗（前端视角）

提醒弹窗运行在独立的 Wails 进程中，前端侧的协作流程：

1. `remindGuard` 检测到当前是提醒进程（`IsRemindWindow()`）时重定向到 `/remind`
2. `pages/remind.vue` 读取样式与模式，从 `remindComponents` 映射中分发对应样式组件，渲染完成后调用 `WindowShow()` 显示窗口，30 秒无操作自动关闭
3. 样式组件约定：接收 `mode` prop（`advance` / `offwork` / `preview`，预览按下班文案展示、不加额外标识），发出 `close` 事件

浏览器中直接调试弹窗样式：

```
http://localhost:5173/remind?style=squidward&mode=offwork
```

新增提醒样式的完整步骤见[根目录 README](../README.md#-扩充新的提醒样式)。
