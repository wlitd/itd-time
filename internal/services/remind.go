package services

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"slices"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// RemindConfig 提醒配置（由前端同步，主进程定时器据此判断提醒时机）
type RemindConfig struct {
	OffTime        string // 下班时间，格式 HH:mm
	AdvanceMinutes int    // 提前提醒分钟数
	RepeatDays     []int  // 重复提醒日（1=周一 ... 7=周日）
	Style          string // 提醒样式
}

// RemindService 提醒服务，管理提醒配置、秒级定时器与弹窗子进程
type RemindService struct {
	ctx      context.Context
	isRemind bool
	style    string
	mode     string
	title    string
	// 以下两个字段从主进程 CLI 传入，供提醒弹窗进程计算文案使用
	offTime        string
	advanceMinutes int

	mu  sync.Mutex
	cfg RemindConfig
}

// NewRemindService 创建提醒服务
func NewRemindService() *RemindService {
	return &RemindService{}
}

// ParseArgs 解析命令行参数，判断是否为提醒弹窗进程
func (s *RemindService) ParseArgs(args []string) {
	for _, arg := range args {
		switch {
		case arg == "--remind":
			s.isRemind = true
		case strings.HasPrefix(arg, "--style="):
			s.style = strings.TrimPrefix(arg, "--style=")
		case strings.HasPrefix(arg, "--mode="):
			s.mode = strings.TrimPrefix(arg, "--mode=")
		case strings.HasPrefix(arg, "--title="):
			s.title = strings.TrimPrefix(arg, "--title=")
		case strings.HasPrefix(arg, "--offtime="):
			s.offTime = strings.TrimPrefix(arg, "--offtime=")
		case strings.HasPrefix(arg, "--advance="):
			s.advanceMinutes, _ = strconv.Atoi(strings.TrimPrefix(arg, "--advance="))
		}
	}
}

// IsRemind 当前进程是否为提醒弹窗进程
func (s *RemindService) IsRemind() bool { return s.isRemind }

// Style 获取提醒弹窗样式
func (s *RemindService) Style() string { return s.style }

// Mode 获取提醒弹窗模式
func (s *RemindService) Mode() string { return s.mode }

// TodoTitle 获取待办标题（仅 todoAdvance 模式有效）
func (s *RemindService) TodoTitle() string { return s.title }

// OffTime 获取下班时间（提醒弹窗进程用于生成文案）
func (s *RemindService) OffTime() string { return s.offTime }

// AdvanceMinutes 获取提前提醒分钟数（提醒弹窗进程用于生成文案）
func (s *RemindService) AdvanceMinutes() int { return s.advanceMinutes }

// Start 启动服务：主进程启动定时器，弹窗进程仅保存 ctx
func (s *RemindService) Start(ctx context.Context) {
	s.ctx = ctx
	if s.isRemind {
		return
	}
	go s.runTicker()
}

// ShowRemind 启动独立的提醒弹窗子进程
// style 为提醒样式（如 squidward），mode 为提醒模式，title 为待办标题（仅 todoAdvance 模式使用）
// 同时将 offTime / advanceMinutes 一并传入，供弹窗进程直接生成文案，无需依赖 store
func (s *RemindService) ShowRemind(style, mode, title string) {
	s.mu.Lock()
	cfg := s.cfg
	s.mu.Unlock()

	exe, err := os.Executable()
	if err != nil {
		return
	}
	args := []string{"--remind", "--style=" + style, "--mode=" + mode}
	if title != "" {
		args = append(args, "--title="+title)
	}
	// 传入 offTime / advanceMinutes，避免弹窗进程依赖 localStorage（进程隔离）
	if cfg.OffTime != "" {
		args = append(args, "--offtime="+cfg.OffTime)
	}
	args = append(args, "--advance="+strconv.Itoa(cfg.AdvanceMinutes))
	_ = exec.Command(exe, args...).Start()
}

// SyncConfig 同步提醒配置（由前端调用）
func (s *RemindService) SyncConfig(offTime string, advanceMinutes int, repeatDays []int, style string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.cfg = RemindConfig{
		OffTime:        offTime,
		AdvanceMinutes: advanceMinutes,
		RepeatDays:     repeatDays,
		Style:          style,
	}
}

// CloseRemind 关闭提醒弹窗（退出当前进程）
func (s *RemindService) CloseRemind() {
	runtime.Quit(s.ctx)
}

// runTicker 主进程内的秒级定时器，在设定时间到达时弹出提醒窗口
func (s *RemindService) runTicker() {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	prev := time.Now()
	for now := range ticker.C {
		s.checkRemind(prev, now)
		prev = now
	}
}

// checkRemind 判断 (prev, now] 区间内是否跨越了提醒时间点
func (s *RemindService) checkRemind(prev, now time.Time) {
	s.mu.Lock()
	cfg := s.cfg
	s.mu.Unlock()

	if cfg.OffTime == "" || len(cfg.RepeatDays) == 0 {
		return
	}

	// 今天是否为重复提醒日（1=周一 ... 7=周日）
	weekday := int(now.Weekday())
	if weekday == 0 {
		weekday = 7
	}
	repeat := false

	repeat = slices.Contains(cfg.RepeatDays, weekday)
	if !repeat {
		return
	}

	target, err := time.ParseInLocation("2006-01-02 15:04",
		fmt.Sprintf("%s %s", now.Format("2006-01-02"), cfg.OffTime), time.Local)
	if err != nil {
		return
	}

	// 提前提醒
	if cfg.AdvanceMinutes > 0 {
		notifyAt := target.Add(-time.Duration(cfg.AdvanceMinutes) * time.Minute)
		if timeCrossed(prev, now, notifyAt) {
			s.ShowRemind(cfg.Style, "advance", "")
		}
	}

	// 下班提醒（倒计时归零）
	if timeCrossed(prev, now, target) {
		s.ShowRemind(cfg.Style, "offwork", "")
	}
}

// timeCrossed 判断时间点 t 是否在 (prev, now] 区间内被跨越
func timeCrossed(prev, now, t time.Time) bool {
	return prev.Before(t) && !now.Before(t)
}
