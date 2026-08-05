package services

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"

	wailsRuntime "github.com/wailsapp/wails/v2/pkg/runtime"
)

// 通过 -ldflags 注入的版本号
var Version = "dev"

type UpdateInfo struct {
	Version     string `json:"version"`
	DownloadUrl string `json:"downloadUrl"`
	ReleaseNote string `json:"releaseNote"`
}

type UpdateResult struct {
	HasUpdate bool        `json:"hasUpdate"`
	Info      *UpdateInfo `json:"info"`
	Error     string      `json:"error"` // 前端可能需要错误信息
}

type UpdaterService struct {
	ctx context.Context
}

func NewUpdaterService() *UpdaterService {
	return &UpdaterService{}
}

// Startup 必须在 app.go 的 startup 中调用，传入 context
func (s *UpdaterService) Startup(ctx context.Context) {
	s.ctx = ctx
}

// CheckForUpdate 检查更新
func (s *UpdaterService) CheckForUpdate() UpdateResult {
	// 1. 配置请求
	client := http.Client{Timeout: 10 * time.Second}
	// 替换为你的真实服务器地址
	url := "https://wlitd.com/updates/latest.json"

	resp, err := client.Get(url)
	if err != nil {
		return UpdateResult{Error: "网络请求失败: " + err.Error()}
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return UpdateResult{Error: "服务器返回错误: " + fmt.Sprintf("%d", resp.StatusCode)}
	}

	// 2. 解析 JSON
	var info UpdateInfo
	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		return UpdateResult{Error: "解析版本信息失败: " + err.Error()}
	}

	// 3. 对比版本
	if info.Version != Version {
		return UpdateResult{HasUpdate: true, Info: &info} // 有新版本
	}

	return UpdateResult{HasUpdate: false} // 无新版本
}

// DownloadAndInstall 下载并安装
func (s *UpdaterService) DownloadAndInstall(url string) error {
	// 1. 创建临时文件
	fileName := filepath.Base(url)
	// 如果是 Windows，假设是 exe。如果是 Mac，假设是 dmg 或 zip
	tempDir := os.TempDir()
	filePath := filepath.Join(tempDir, fileName)

	// 2. 发起下载请求（使用带超时的客户端）
	dlClient := http.Client{Timeout: 30 * time.Minute}
	resp, err := dlClient.Get(url)
	if err != nil {
		return fmt.Errorf("下载失败: %v", err)
	}
	defer resp.Body.Close()

	// 3. 创建本地文件
	out, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("创建文件失败: %v", err)
	}
	defer out.Close()

	// 4. 带进度的拷贝
	// 获取文件总大小用于计算进度
	size := resp.ContentLength
	counter := &writeCounter{
		ctx:      s.ctx,
		total:    size,
		lastEmit: 0,
	}

	// 执行拷贝
	if _, err = io.Copy(out, io.TeeReader(resp.Body, counter)); err != nil {
		return fmt.Errorf("文件写入失败: %v", err)
	}

	// 5. 运行安装包
	return s.runInstaller(filePath)
}

// runInstaller 运行安装程序
func (s *UpdaterService) runInstaller(filePath string) error {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "windows":
		// Windows: 直接启动安装包
		cmd = exec.Command(filePath)
	case "darwin":
		// Mac: 打开 dmg 镜像，用户手动拖拽
		cmd = exec.Command("open", filePath)
	case "linux":
		cmd = exec.Command("xdg-open", filePath)
	default:
		return fmt.Errorf("不支持的操作系统")
	}

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("启动安装程序失败: %v", err)
	}

	return nil
}

// ================= 辅助结构体：用于计算下载进度 =================

type writeCounter struct {
	ctx      context.Context
	total    int64
	written  int64
	lastEmit int64 // 用来节流，避免发送太多事件导致 UI 卡顿
}

func (wc *writeCounter) Write(p []byte) (int, error) {
	n := len(p)
	wc.written += int64(n)

	// 每下载 1% 或者每 100KB 发送一次进度，避免过于频繁
	if wc.total > 0 {
		currentPercent := float64(wc.written) / float64(wc.total) * 100
		if int(currentPercent) > int(wc.lastEmit) {
			wailsRuntime.EventsEmit(wc.ctx, "update-download-progress", int(currentPercent))
			wc.lastEmit = int64(int(currentPercent))
		}
	}
	return n, nil
}
