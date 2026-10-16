// Package api 汇聚所有对前端暴露的 Wails Service（绑定面）。
// Service 是薄封装：只做参数校验、调用 domain/storage、并触发重排等编排，
// 不含业务规则本身（规则在 domain），也不含持久化细节（在 storage）。
package api

import (
	"itdtime/internal/domain"
	"itdtime/internal/storage"
)

// SettingsService 暴露配置与主题读写。
type SettingsService struct {
	store *storage.Store
}

// NewSettingsService 构造 SettingsService。
func NewSettingsService(store *storage.Store) *SettingsService {
	return &SettingsService{store: store}
}

// GetAppData 返回当前全量持久化数据快照，供前端首屏一次性加载。
func (s *SettingsService) GetAppData() domain.AppData {
	return s.store.Data()
}

// SaveOffWorkConfig 更新下班配置（启用开关/下班时间/重复日）并落盘。
// 真正的「重算下次提醒时间」在 reminder.Coordinator 接入后由此处触发（M4）。
func (s *SettingsService) SaveOffWorkConfig(cfg domain.OffWorkConfig) error {
	return s.store.Update(func(data *domain.AppData) error {
		data.OffWork = cfg
		return nil
	})
}

// SaveRemindSettings 更新全局提醒偏好（提前量 + 样式）并落盘。
func (s *SettingsService) SaveRemindSettings(cfg domain.RemindSettings) error {
	return s.store.Update(func(data *domain.AppData) error {
		data.Remind = cfg
		return nil
	})
}

// SaveTheme 更新主题/语言偏好并落盘（跨窗口单一事实源）。
func (s *SettingsService) SaveTheme(theme domain.ThemePrefs) error {
	return s.store.Update(func(data *domain.AppData) error {
		data.Theme = theme
		return nil
	})
}

// SaveUpdatePrefs 更新更新相关偏好（启动检查开关）并落盘。
func (s *SettingsService) SaveUpdatePrefs(cfg domain.UpdatePrefs) error {
	return s.store.Update(func(data *domain.AppData) error {
		data.Update = cfg
		return nil
	})
}
