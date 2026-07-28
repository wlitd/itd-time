package main

import (
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// remindConfig 前端同步过来的提醒配置
type remindConfig struct {
	offTime        string // 下班时间，格式 HH:mm
	advanceMinutes int    // 提前提醒分钟数
	repeatDays     []int  // 重复提醒日（1=周一 ... 7=周日）
	style          string // 提醒样式
}

// ShowRemind 启动一个独立的提醒弹窗进程
// style 为提醒样式（如 squidward），mode 为提醒模式（advance/offwork/preview）
func (a *App) ShowRemind(style string, mode string) {
	exe, err := os.Executable()
	if err != nil {
		return
	}
	_ = exec.Command(exe, "--remind", "--style="+style, "--mode="+mode).Start()
}

// SyncRemindConfig 由前端同步提醒配置，供主进程定时器判断提醒时机
func (a *App) SyncRemindConfig(offTime string, advanceMinutes int, repeatDays []int, style string) {
	a.remindMu.Lock()
	defer a.remindMu.Unlock()
	a.remindCfg = remindConfig{
		offTime:        offTime,
		advanceMinutes: advanceMinutes,
		repeatDays:     repeatDays,
		style:          style,
	}
}

// IsRemindWindow 当前进程是否为提醒弹窗进程
func (a *App) IsRemindWindow() bool {
	return a.isRemind
}

// GetRemindStyle 获取提醒弹窗的样式
func (a *App) GetRemindStyle() string {
	return a.remindStyle
}

// GetRemindMode 获取提醒弹窗的模式
func (a *App) GetRemindMode() string {
	return a.remindMode
}

// CloseRemind 关闭提醒弹窗（退出提醒进程）
func (a *App) CloseRemind() {
	runtime.Quit(a.ctx)
}

// runRemindTicker 主进程内的秒级定时器，在设定时间到达时弹出提醒窗口
func (a *App) runRemindTicker() {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	prev := time.Now()
	for now := range ticker.C {
		a.checkRemind(prev, now)
		prev = now
	}
}

// checkRemind 判断 (prev, now] 区间内是否跨越了提醒时间点
func (a *App) checkRemind(prev, now time.Time) {
	a.remindMu.Lock()
	cfg := a.remindCfg
	a.remindMu.Unlock()

	if cfg.offTime == "" || len(cfg.repeatDays) == 0 {
		return
	}

	// 今天是否为重复提醒日（1=周一 ... 7=周日）
	weekday := int(now.Weekday())
	if weekday == 0 {
		weekday = 7
	}
	repeat := false
	for _, d := range cfg.repeatDays {
		if d == weekday {
			repeat = true
			break
		}
	}
	if !repeat {
		return
	}

	target, err := time.ParseInLocation("2006-01-02 15:04",
		fmt.Sprintf("%s %s", now.Format("2006-01-02"), cfg.offTime), time.Local)
	if err != nil {
		return
	}

	// 提前提醒
	if cfg.advanceMinutes > 0 {
		notifyAt := target.Add(-time.Duration(cfg.advanceMinutes) * time.Minute)
		if timeCrossed(prev, now, notifyAt) {
			a.ShowRemind(cfg.style, "advance")
		}
	}

	// 下班提醒（倒计时归零）
	if timeCrossed(prev, now, target) {
		a.ShowRemind(cfg.style, "offwork")
	}
}

// timeCrossed 判断时间点 t 是否在 (prev, now] 区间内被跨越
func timeCrossed(prev, now, t time.Time) bool {
	return prev.Before(t) && !now.Before(t)
}
