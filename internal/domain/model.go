// Package domain 定义应用的核心数据模型与调度规则。
// 本包不依赖 Wails / 存储 / 窗口等任何基础设施，保证逻辑可被纯单元测试覆盖。
package domain

import "time"

// RemindStyle 提醒样式标识，可插拔。
type RemindStyle string

const (
	// RemindStyleSquidward 海绵宝宝（痞老板）风格提醒卡片。
	RemindStyleSquidward RemindStyle = "squidward"
	// RemindStyleMinimal 简约风格提醒卡片。
	RemindStyleMinimal RemindStyle = "minimal"
)

// RemindKind 提醒类型。
type RemindKind string

const (
	// RemindKindOffWork 下班提醒。
	RemindKindOffWork RemindKind = "offwork"
	// RemindKindOffWorkAdvance 下班提前提醒。
	RemindKindOffWorkAdvance RemindKind = "offworkAdvance"
	// RemindKindTodo 代办提醒（在截止时间前 AdvanceMinutes 触发；为 0 时即到期时刻）。
	RemindKindTodo RemindKind = "todo"
)

// RemindSettings 提醒的通用偏好，全局生效——同时作用于下班提醒与代办提醒。
// 之所以独立于 OffWorkConfig：提前量与卡片样式并非「下班」专属，
// 代办也复用同一份提前量，避免把全局规则耦合进某一类配置里。
type RemindSettings struct {
	// AdvanceMinutes 提前提醒分钟数，0 表示不提前。
	AdvanceMinutes int `json:"advanceMinutes"`
	// Style 提醒卡片样式。
	Style RemindStyle `json:"style"`
}

// OffWorkConfig 下班提醒配置——只描述「什么时候算下班」。
type OffWorkConfig struct {
	// Enabled 是否启用下班提醒。
	Enabled bool `json:"enabled"`
	// OffTime 下班时间，格式 "15:04"。
	OffTime string `json:"offTime"`
	// RepeatDays 重复日，取值 0=周日 ... 6=周六；空切片表示每天。
	RepeatDays []time.Weekday `json:"repeatDays"`
}

// Todo 代办事项。
type Todo struct {
	// ID 唯一标识。
	ID string `json:"id"`
	// Title 标题。
	Title string `json:"title"`
	// Deadline 截止时间（Unix 毫秒）。
	Deadline int64 `json:"deadline"`
	// Completed 是否已完成。
	Completed bool `json:"completed"`
	// Reminded 是否已触发过提醒（由后端回写，防止重启后重复弹窗）。
	Reminded bool `json:"reminded"`
}

// ThemePrefs 主题与语言偏好。多窗口共享，作为跨窗口的单一事实源。
type ThemePrefs struct {
	// Dark 是否为深色模式。
	Dark bool `json:"dark"`
	// Locale 语言代码，如 "zh-CN" / "en-US"。
	Locale string `json:"locale"`
}

// UpdatePrefs 更新相关偏好。
type UpdatePrefs struct {
	// CheckOnLaunch 应用启动时是否静默检查更新（仅发现新版本时打扰）。
	CheckOnLaunch bool `json:"checkOnLaunch"`
}

// AppData 应用持久化数据的根结构，对应存储中的单个 JSON 文件。
type AppData struct {
	// SchemaVersion 数据结构版本，用于后续迁移。
	SchemaVersion int `json:"schemaVersion"`
	// OffWork 下班提醒配置。
	OffWork OffWorkConfig `json:"offWork"`
	// Remind 全局提醒偏好（提前量 + 样式）。
	Remind RemindSettings `json:"remind"`
	// Todos 代办列表。
	Todos []Todo `json:"todos"`
	// Theme 主题偏好。
	Theme ThemePrefs `json:"theme"`
	// Update 更新相关偏好（启动检查等）。
	Update UpdatePrefs `json:"update"`
	// SkippedVersion 被用户跳过、不再提醒更新的版本号。
	SkippedVersion string `json:"skippedVersion"`
}

// DefaultAppData 返回一份合理默认值，用于首次启动或数据缺失时。
func DefaultAppData() AppData {
	return AppData{
		SchemaVersion: 1,
		OffWork: OffWorkConfig{
			Enabled:    true,
			OffTime:    "17:30",
			RepeatDays: []time.Weekday{time.Monday, time.Tuesday, time.Wednesday, time.Thursday, time.Friday},
		},
		Remind: RemindSettings{
			AdvanceMinutes: 10,
			Style:          RemindStyleSquidward,
		},
		Todos:          []Todo{},
		Theme:          ThemePrefs{Dark: false, Locale: "zh-CN"},
		Update:         UpdatePrefs{CheckOnLaunch: true},
		SkippedVersion: "",
	}
}
