# ItdTime 前端（frontend/）

Vue 3.5 + TypeScript + Vite 8 + PrimeVue 4 + UnoCSS，运行于 Wails v3 的 WebView 中。

## 脚本

| 命令 | 说明 |
| --- | --- |
| `pnpm dev` | 纯浏览器模式（Vite Dev Server，端口 9245），Wails 绑定自动降级 |
| `pnpm build:dev` | 类型检查 + 不压缩构建（供 `wails3 dev` 内嵌使用） |
| `pnpm build` | 类型检查 + 生产构建（`vue-tsc && vite build`） |
| `pnpm preview` | 预览构建产物 |

> 桌面端日常开发请在**仓库根目录**跑 `wails3 task dev`（会代管本目录的构建与绑定生成）。

## 目录结构

```
src/
├── core/                # 基础设施层
│   ├── i18n/            # vue-i18n 实例 + 集中式 messages（zh-CN / en）
│   ├── theme/           # PrimeVue 自定义主题 preset（DESIGN.md 的代码化）
│   ├── bindings/        # 对生成绑定的统一 re-export（见「与 Wails 的协作」）
│   └── stores/          # 全局 Pinia store（薄缓存，数据主权在 Go 端）
├── features/            # 按领域切分的业务组件，就近放 types.ts
│   ├── countdown/       # LED 倒计时卡片
│   ├── remind/          # 提醒窗：RemindWindow 分发 + registry 样式注册表
│   ├── settings/        # 设置弹窗（含 RepeatSetting 子组件）
│   ├── shell/           # 自绘标题栏 AppHeader / ThemeToggle
│   ├── todos/           # 待办列表（TodoList / TodoItem）
│   └── update/          # 应用内自更新（UpdateDialog + useUpdate）
├── pages/MainView.vue   # 主页面
├── shared/              # env.ts 运行环境判定、datetime.ts 日期工具
└── styles/main.css      # CSS 变量 + LED 字体 + 全局修正
```

## 与 Wails 的协作

- **绑定生成**：`wails3 dev/build` 会生成 `frontend/bindings/`（TS 绑定，已 gitignore，勿手改）。业务代码统一从 `src/core/bindings` 导入，不直接引用生成目录。
- **环境降级**：`shared/env.ts` 的 `isWails()` 判定是否运行在桌面宿主；纯浏览器（`pnpm dev`）下绑定调用降级为占位/内存态，方便 UI 调试。
- **提醒窗复用同一构建**：Go 端 Coordinator 以 `/?remind=1&kind=…&id=…` 打开独立小窗，`App.vue` 据 `isRemindWindow()` 分发渲染 `RemindWindow`。浏览器里可直接访问该 URL 预览提醒卡片。

## 自动导入

- `unplugin-auto-import`：vue / vue-router / vueuse / vue-i18n 常用 API 免 import（`ref`、`computed`、`useToast` 等）。
- `unplugin-vue-components`：PrimeVue 组件按需自动注册。
- 声明文件 `src/auto-imports.d.ts`、`src/components.d.ts` 已入库（供 `vue-tsc` 在 Vite 生成前完成类型检查），随构建自动更新。
- ⚠️ `import type` **不受**自动导入管理，类型仍需显式导入。

## 编码规范（vue3-ts-coding-style）

- 组件拆为 `Xxx.vue` + 就近 `types.ts`；Props 用 **interface** 定义并经 `defineProps<XxxProps>()` **解构**（默认值写在解构处）
- Emits 用**具名元组**类型声明：`defineEmits<{ change: [files: UploadFileItem[]] }>()`
- 函数一律 `function` 声明 + **显式返回类型**（无返回标 `void`）；箭头函数仅限内联回调
- `ref` / `computed` 显式标注泛型类型；组件对外导出实例类型 `type XxxInst = InstanceType<typeof Xxx>`
- 不使用分号；不留未使用的变量与导入

## 设计系统

[DESIGN.md](./DESIGN.md) 是 UI 的**唯一事实源**（色板、字体、间距、组件约定、深浅主题）。代码侧对应实现：

| 事实源条目 | 实现位置 |
| --- | --- |
| PrimeVue 主题令牌 | `src/core/theme/preset.ts` |
| 工具类 / 断点 / 图标 | `uno.config.ts` |
| CSS 变量 / LED 字体 | `src/styles/main.css` |

## 国际化

- 文案集中在 `src/core/i18n/messages.ts`（`zh-CN` / `en` 双语对照）
- 语言偏好由 Go 端持久化（`ThemePrefs.locale`），启动后经 store 水合覆盖，切换即时生效
