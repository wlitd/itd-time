package domain

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"
)

// Fire 是一个待触发的提醒时刻及其上下文。
// 调度器只负责「算出未来会发生哪些提醒」，不负责「何时真正弹出」——
// 后者由 reminder.Coordinator 用定时器装载（M4）。
type Fire struct {
	// At 触发的绝对时刻（本地时区）。
	At time.Time
	// Kind 提醒类型。
	Kind RemindKind
	// TodoID 仅代办类提醒有效，标识对应的代办。
	TodoID string
}

// ParseClock 解析 "HH:MM" 为一天内的分钟数（自 00:00 起）。
func ParseClock(s string) (int, error) {
	parts := strings.Split(strings.TrimSpace(s), ":")
	if len(parts) != 2 {
		return 0, fmt.Errorf("invalid clock %q, want HH:MM", s)
	}
	h, err := strconv.Atoi(strings.TrimSpace(parts[0]))
	if err != nil || h < 0 || h > 23 {
		return 0, fmt.Errorf("invalid hour in clock %q", s)
	}
	m, err := strconv.Atoi(strings.TrimSpace(parts[1]))
	if err != nil || m < 0 || m > 59 {
		return 0, fmt.Errorf("invalid minute in clock %q", s)
	}
	return h*60 + m, nil
}

// matchesRepeat 判断给定星期是否命中重复规则；空规则视为每天。
func matchesRepeat(day time.Weekday, repeat []time.Weekday) bool {
	if len(repeat) == 0 {
		return true
	}
	for _, d := range repeat {
		if d == day {
			return true
		}
	}
	return false
}

// NextFires 计算 [now, now+horizon) 区间内所有应触发的提醒，按时间升序返回。
// 纯函数：不读时钟、不写状态，结果仅由 data 与 now 决定。
func NextFires(data AppData, now time.Time, horizon time.Duration) []Fire {
	advance := data.Remind.AdvanceMinutes
	var fires []Fire
	fires = append(fires, offWorkFires(data.OffWork, advance, now, horizon)...)
	fires = append(fires, todoFires(data.Todos, advance, now, horizon)...)
	sort.Slice(fires, func(i, j int) bool {
		if !fires[i].At.Equal(fires[j].At) {
			return fires[i].At.Before(fires[j].At)
		}
		return fires[i].Kind < fires[j].Kind
	})
	return fires
}

// offWorkFires 计算下班主提醒与提前提醒。
func offWorkFires(cfg OffWorkConfig, advance int, now time.Time, horizon time.Duration) []Fire {
	if !cfg.Enabled {
		return nil
	}
	minutes, err := ParseClock(cfg.OffTime)
	if err != nil {
		return nil
	}
	end := now.Add(horizon)
	var out []Fire
	// 逐日扫描 horizon 覆盖到的每一天（含今天）。
	for d := 0; ; d++ {
		day := time.Date(now.Year(), now.Month(), now.Day()+d, 0, 0, 0, 0, now.Location())
		if day.After(end) {
			break
		}
		if !matchesRepeat(day.Weekday(), cfg.RepeatDays) {
			continue
		}
		offAt := day.Add(time.Duration(minutes) * time.Minute)
		if offAt.After(now) && !offAt.After(end) {
			out = append(out, Fire{At: offAt, Kind: RemindKindOffWork})
		}
		if advance > 0 {
			advAt := offAt.Add(-time.Duration(advance) * time.Minute)
			if advAt.After(now) && !advAt.After(end) {
				out = append(out, Fire{At: advAt, Kind: RemindKindOffWorkAdvance})
			}
		}
	}
	return out
}

// todoFires 计算未完成、未提醒过、且提醒时刻落在区间内的代办。
// 与下班提醒共用全局提前量：代办在 (deadline - advance) 触发一次；
// advance 为 0 时即到期那一刻。已提醒（Reminded）的代办不再重复弹出。
func todoFires(todos []Todo, advance int, now time.Time, horizon time.Duration) []Fire {
	end := now.Add(horizon)
	advanceDur := time.Duration(advance) * time.Minute
	var out []Fire
	for _, t := range todos {
		if t.Completed || t.Reminded {
			continue
		}
		at := time.UnixMilli(t.Deadline).In(now.Location()).Add(-advanceDur)
		if at.After(now) && !at.After(end) {
			out = append(out, Fire{At: at, Kind: RemindKindTodo, TodoID: t.ID})
		}
	}
	return out
}
