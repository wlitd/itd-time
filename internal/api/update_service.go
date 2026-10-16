package api

import (
	"context"
	"os"
	"path/filepath"

	"github.com/wailsapp/wails/v3/pkg/application"
	"github.com/wailsapp/wails/v3/pkg/updater"
	"github.com/wailsapp/wails/v3/pkg/updater/providers/endpoint"

	"itdtime/internal/domain"
	"itdtime/internal/storage"
)

// updateManifestURL 指向发布在云服务器上的 Wails Update Manifest（静态 JSON）。
// 采用 V3 原生 updater + endpoint provider：
//   - manifest 描述最新版本及其各平台产物，产物附带 sha256 摘要，由 updater 流式校验；
//   - 安装动作是「就地二进制替换」而非运行安装器，故要求以 user scope 安装
//     （$LOCALAPPDATA\Programs\ItdTime），进程对安装目录有写权限，免 UAC 提权。
const updateManifestURL = "https://wlitd.com/updates/manifest.json"

// updateChannel 发布渠道；endpoint provider 会附带 channel 查询参数并在客户端校验。
const updateChannel = "stable"

// UpdateInfo 版本检查结果，供前端展示。仅暴露 UI 需要的字段，
// 避免把框架内部 updater.Release 的完整结构泄漏到绑定面。
type UpdateInfo struct {
	// HasUpdate 是否存在可更新版本（已跳过或已最新时为 false）。
	HasUpdate bool `json:"hasUpdate"`
	// CurrentVersion 当前运行版本。
	CurrentVersion string `json:"currentVersion"`
	// Version 可获取的版本号。
	Version string `json:"version"`
	// Name 版本展示名。
	Name string `json:"name"`
	// Notes 更新说明。
	Notes string `json:"notes"`
	// Size 产物字节数（未知时为 0）。
	Size int64 `json:"size"`
}

// UpdateService 暴露自更新能力：检查、下载安装、重启应用、跳过版本。
// 下载进度等细粒度状态由前端订阅 wails:updater:* 事件获取，不在此重复轮询。
type UpdateService struct {
	store          *storage.Store
	app            *application.App
	currentVersion string
}

// NewUpdateService 构造 UpdateService。currentVersion 由装配层经 -ldflags 注入。
func NewUpdateService(store *storage.Store, currentVersion string) *UpdateService {
	return &UpdateService{store: store, currentVersion: currentVersion}
}

// ServiceStartup 实现 application.ServiceStartup：初始化 updater 并恢复已跳过版本。
func (s *UpdateService) ServiceStartup(_ context.Context, _ application.ServiceOptions) error {
	s.app = application.Get()

	// 清理上一次就地更新遗留的旧二进制：Windows 下 updater 把正在运行的目标 exe
	// 改名为 "<exe>.old.<随机>" 后再落新文件，改名后的旧文件当时仍被占用无法删除，
	// 会残留到下次启动。此时新进程已接管、旧文件不再被锁，可安全清扫。
	cleanupStaleBinaries()

	provider, err := endpoint.New(endpoint.Config{
		URL:     updateManifestURL,
		Channel: updateChannel,
	})
	if err != nil {
		return err
	}

	if err := s.app.Updater.Init(updater.Config{
		CurrentVersion: s.currentVersion,
		Providers:      []updater.Provider{provider},
	}); err != nil {
		return err
	}

	// 跨进程恢复「跳过此版本」：updater 的 skip 仅存于内存，重启即失效，
	// 故用持久化的 SkippedVersion 回填，保持与 V2 一致的行为。
	if skipped := s.store.Data().SkippedVersion; skipped != "" {
		s.app.Updater.SkipVersion(skipped)
	}
	return nil
}

// cleanupStaleBinaries 删除与当前可执行文件同目录、形如 "<exe>.old.*" 的更新残留。
// 尽力而为：任何错误都忽略（文件可能仍被占用），留待下次启动再清。
func cleanupStaleBinaries() {
	exePath, err := os.Executable()
	if err != nil {
		return
	}
	stale, err := filepath.Glob(exePath + ".old.*")
	if err != nil {
		return
	}
	for _, path := range stale {
		_ = os.Remove(path)
	}
}

// CheckForUpdate 主动检查更新。返回 HasUpdate=false 表示已最新或该版本已被跳过。
func (s *UpdateService) CheckForUpdate() (UpdateInfo, error) {
	rel, err := s.app.Updater.Check(context.Background())
	if err != nil {
		return UpdateInfo{}, err
	}
	if rel == nil {
		return UpdateInfo{HasUpdate: false, CurrentVersion: s.currentVersion}, nil
	}
	return UpdateInfo{
		HasUpdate:      true,
		CurrentVersion: s.currentVersion,
		Version:        rel.Version,
		Name:           rel.Name,
		Notes:          rel.Notes,
		Size:           rel.Artifact.Size,
	}, nil
}

// DownloadAndInstall 下载上一步 Check 命中的版本，校验 sha256 并暂存待替换。
// 完成后进入 Ready 态，需再调用 ApplyUpdate 才会真正替换二进制并重启。
func (s *UpdateService) DownloadAndInstall() error {
	return s.app.Updater.DownloadAndInstall(context.Background())
}

// ApplyUpdate 用已暂存的新版本替换当前二进制并重启（会退出当前进程）。
func (s *UpdateService) ApplyUpdate() error {
	return s.app.Updater.Restart(context.Background())
}

// SkipVersion 跳过指定版本：不再提醒该版本更新，并持久化以跨重启生效。
func (s *UpdateService) SkipVersion(version string) error {
	s.app.Updater.SkipVersion(version)
	return s.store.Update(func(data *domain.AppData) error {
		data.SkippedVersion = version
		return nil
	})
}

// CurrentVersion 返回当前运行版本，供「关于」页展示。
func (s *UpdateService) CurrentVersion() string {
	return s.currentVersion
}
