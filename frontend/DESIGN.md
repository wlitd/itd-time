# DESIGN.md — ItdTime 设计系统（唯一事实源）

> 桌面倒计时 / 下班提醒工具的设计规范。技术栈：**Vue 3.5 + PrimeVue 4 + UnoCSS**。
> 所有色板、字体、间距、组件约定以本文件为准；代码侧对应实现见
> `src/core/theme/preset.ts`（PrimeVue）、`uno.config.ts`（工具类）、`src/styles/main.css`（CSS 变量）。

---

## 1. 设计基调

- **人格**：克制、专注、带一点幽默（章鱼哥 / “我免费啦”）。视觉安静，动效轻微，信息层级清晰。
- **形态**：紧凑型桌面窗口（主窗约 `820 × 600`），无边框 + 自绘头部拖拽区；提醒为右下角透明置顶小卡。
- **双主题**：light / dark 一视同仁，均保证正文对比度 ≥ 4.5:1。主题由 `<html>` 上的 `.dark` 类驱动，跨窗口以 Go 端 `ThemePrefs` 为单一事实源。

---

## 2. 色彩

### 2.1 品牌主色 —— 玫红 Rose

> 背景渐变共享粉/紫/暖橙三条色带；primary 锚在**品红侧的玫红**（~330°，Tailwind pink 阶），
> 既贴合背景的粉带（`#EF4E7B` / `#E96F92`），又与 danger 纯红（~0°）拉开距离、避免歧义。

| Token | Light 用途 | Dark 用途 |
|-------|-----------|-----------|
| primary-600 `#DB2777` | 主按钮 / 强调文字 | — |
| primary-400 `#F472B6` | — | 主按钮 / 强调文字 |
| primary-50 `#FDF2F8` | 弱背景 / hover 底色 | — |
| primary-950 `#500724` | — | 深色主色底 |

完整 50→950 阶梯见 `preset.ts` 的 `rose`；`uno.config.ts` 的 `theme.colors.primary` 与之一致。
PrimeVue 的 `semantic.primary` 与 `colorScheme.{light,dark}.primary` 已按上表映射。

- **辅助色 accent**：橙 `#F97316`，仅用于「超时 / 紧急 / 下班倒计时归零」等警示，克制使用。
- **语义危险 danger**：红 `#DC2626`（light）/ `#F87171`（dark），与玫红主色以色相 + 图标/文案双重区分，不单靠颜色。
- **LED 数字**：主色发光 `.led-glow`（`text-primary-400` + 粉光 text-shadow），深色下更醒目。

### 2.2 中性表面（CSS 变量，随主题切换）

| 变量 | Light | Dark | 用途 |
|------|-------|------|------|
| `--app-bg` | `#FFFFFF`（白底） | `#1B2947`（深藏青） | 渐变下方的兜底纯色 |
| `--app-bg-image` | 8 段低透明度横向彩虹渐变（90deg，各色 0.1 alpha 叠白） | `#1B2947 → #75517D → #E96F92 → #F7F7B6` 垂直渐变 | 页面背景（body） |
| `--app-surface` | `#FFFFFF` | `#131D21` | 卡片 / 面板 |
| `--app-text` | `#0F2E2B` | `#E7F0EE` | 正文 |
| `--app-muted` | `#5F7178` | `#8BA0A6` | 次要文字 |
| `--app-border` | `#D9E6E2` | `#24333A` | 描边 / 分割 |

背景由 `--app-bg`（兜底色）+ `--app-bg-image`（主题渐变）叠加而成，在 `body` 上以 `background-color` + `background-image` 呈现。
卡片（`--app-surface`）为不透明面板，浮于渐变之上；交互/正文文字均落在卡片或头部深色区，保证对比度。
提醒窗以 `background: transparent !important` 覆盖，不渲染该渐变。

PrimeVue 组件沿用 Aura 的 `surface` 刻度；自绘区域用上述 CSS 变量，保证与组件库协调。

---

## 3. 字体

| 角色 | 字体栈 | 备注 |
|------|--------|------|
| 正文 / UI | `'Plus Jakarta Sans', system-ui, …, 'PingFang SC', 'Microsoft YaHei', sans-serif` | 中英混排回退系统中文栈 |
| 倒计时数字 | `'LED', ui-monospace, monospace`（`.led`） | 本地内嵌 `assets/font/LED.ttf`，`font-variant-numeric: tabular-nums` 防抖动 |

字号（工具类）：页面标题 `text-2xl/3xl`；卡片标题 `text-lg`；正文 `text-sm`；辅助 `text-xs`。

---

## 4. 间距与圆角

- **间距**：基于 4px 栅格，UnoCSS 默认刻度。卡片内边距 `p-4`～`p-6`；元素间隙 `gap-2/3/4`。
- **圆角**：卡片 `rounded-2xl`（16px）；按钮/输入沿用 PrimeVue 默认 `rounded-md`；提醒小卡 `rounded-xl`。
- **阴影**：light 用 `shadow-sm`/`shadow-lg`（悬浮层）；dark 主要靠描边与背景对比，弱化阴影。
- **可点击反馈**：一律 `cursor-pointer` + `transition-colors duration-200`；hover 改背景/边框色，不用位移式 scale。

---

## 5. 布局结构（主窗）

```
┌───────────────────────────────────────────┐
│ AppHeader（frameless 拖拽区 · wails-drag）  │  ← 标题 / 主题 / 语言 / 更新 / 自启 / 菜单
├───────────────────────────────────────────┤
│  ┌─────────────────┐  ┌─────────────────┐  │
│  │  CountdownCard   │  │    TodoList     │  │
│  │  LED 倒计时 + 设置│  │  CRUD / 逾期    │  │
│  └─────────────────┘  └─────────────────┘  │
└───────────────────────────────────────────┘
```

- 头部按钮等交互区标记 `wails-no-drag`，其余头部 `wails-drag` 可拖动窗口。
- 主体两栏：左倒计时卡，右代办列表；窄屏（<720px）纵向堆叠。

---

## 6. 提醒窗（透明置顶）

- 尺寸 `300 × 120`，屏幕右下角，`BackgroundTypeTransparent`，无边框、置顶、不可缩放。
- 可插拔样式（`RemindStyle`）经 `features/remind/registry.ts` 映射到组件：
  - **squidward**：章鱼哥形象卡（`SquidwardRemind.vue`），标题/描述 + 「朕知道了 / 我就硬拖」。
  - **minimal**：克制纯文字卡（`MinimalRemind.vue`）。
- 卡体半透明表面 + 圆角 + 明显描边（保证 light 下也可见）。文案由 `kind` + i18n 模板重建。

---

## 7. 组件约定（PrimeVue 4）

| 场景 | 组件 |
|------|------|
| 主/次按钮、图标按钮 | `Button`（`severity` / `text` / `rounded`） |
| 设置面板 / 更新弹窗 | `Dialog`（`modal`，头部含关闭） |
| 下班时间 | `DatePicker`（`timeOnly` + `hourFormat="24"`） |
| 提前分钟 | `InputNumber`（`useGrouping=false`，`min=0`） |
| 代办截止时间 | `DatePicker`（`showTime`，`hourFormat=24`） |
| 重复日 | 自绘星期 chip（`RepeatSetting.vue`） |
| 提醒样式 | 自绘卡片单选（`RemindStyleSelect.vue`） |
| 开关（自启/启用下班） | `ToggleSwitch` |
| 菜单 / 语言 / 更多操作 | `Menu` + `Button(text)` |
| 全局提示 | `Toast` |
| 更新进度 | `ProgressBar` |

- 图标统一 **Lucide**（`i-lucide:*`），语言用 **flag**（`i-flag:*`）；禁止用 emoji 当图标。
- PrimeVue 组件**按需显式 import**（自动注册仅覆盖 `src` 下自研组件）。

---

## 8. 动效

- 过渡 150–300ms，`ease-out`；主题切换不加特殊动画，仅颜色过渡。
- 尊重 `prefers-reduced-motion`：倒计时数字变化不做弹跳，仅等宽替换。
- 提醒卡出现：轻微淡入 + 位移 8px；关闭交由后端 `Dismiss`。

---

## 9. 反模式（避免）

- ❌ 混用多种圆角/阴影/主色阶；✅ 一律走 token。
- ❌ light 模式用过透的半透明表面（<80% 不透明度）导致文字的对比不足。
- ❌ 仅用颜色区分状态（逾期）；✅ 同时给出文案/图标。
- ❌ 在头部按钮上遗留 `wails-drag`，导致按钮点不动 —— 交互元素必须 `wails-no-drag`。
