package domain

import (
	"testing"
	"time"
)

var cst = time.FixedZone("CST", 8*3600)

func at(y int, mo time.Month, d, h, mi int) time.Time {
	return time.Date(y, mo, d, h, mi, 0, 0, cst)
}

func TestParseClock(t *testing.T) {
	cases := []struct {
		in      string
		wantMin int
		wantErr bool
	}{
		{"17:30", 17*60 + 30, false},
		{"00:00", 0, false},
		{"23:59", 23*60 + 59, false},
		{" 9:05 ", 9*60 + 5, false},
		{"24:00", 0, true},
		{"12:60", 0, true},
		{"17", 0, true},
		{"abc:def", 0, true},
	}
	for _, c := range cases {
		got, err := ParseClock(c.in)
		if c.wantErr {
			if err == nil {
				t.Errorf("ParseClock(%q) expected error, got %d", c.in, got)
			}
			continue
		}
		if err != nil {
			t.Errorf("ParseClock(%q) unexpected error: %v", c.in, err)
		}
		if got != c.wantMin {
			t.Errorf("ParseClock(%q) = %d, want %d", c.in, got, c.wantMin)
		}
	}
}

func baseData() AppData {
	d := DefaultAppData()
	d.OffWork = OffWorkConfig{
		Enabled:    true,
		OffTime:    "17:30",
		RepeatDays: []time.Weekday{time.Monday, time.Tuesday, time.Wednesday, time.Thursday, time.Friday},
	}
	d.Remind = RemindSettings{AdvanceMinutes: 0, Style: RemindStyleSquidward}
	d.Todos = nil
	return d
}

func TestOffWorkFires_WeekdaySameDay(t *testing.T) {
	data := baseData()
	now := at(2026, time.October, 14, 10, 0) // 周三 10:00
	fires := NextFires(data, now, 24*time.Hour)
	if len(fires) != 1 {
		t.Fatalf("want 1 fire, got %d: %+v", len(fires), fires)
	}
	want := at(2026, time.October, 14, 17, 30)
	if !fires[0].At.Equal(want) || fires[0].Kind != RemindKindOffWork {
		t.Errorf("got %+v, want offwork@%v", fires[0], want)
	}
}

func TestOffWorkFires_AlreadyPassedToday_SkipsToNext(t *testing.T) {
	data := baseData()
	now := at(2026, time.October, 14, 18, 0) // 周三 18:00，已过下班点
	fires := NextFires(data, now, 24*time.Hour)
	// 今天已过；周四(15) 17:30 应命中。
	if len(fires) != 1 {
		t.Fatalf("want 1 fire, got %d: %+v", len(fires), fires)
	}
	want := at(2026, time.October, 15, 17, 30)
	if !fires[0].At.Equal(want) {
		t.Errorf("got %v, want %v", fires[0].At, want)
	}
}

func TestOffWorkFires_WeekendExcluded(t *testing.T) {
	data := baseData()
	now := at(2026, time.October, 16, 10, 0) // 周五 10:00
	// 未来 72h 覆盖周五到周一凌晨；周末不在重复日，只应有周五那一次。
	fires := NextFires(data, now, 72*time.Hour)
	for _, f := range fires {
		if f.At.Weekday() == time.Saturday || f.At.Weekday() == time.Sunday {
			t.Errorf("weekend fire leaked: %+v", f)
		}
	}
	if len(fires) != 1 || !fires[0].At.Equal(at(2026, time.October, 16, 17, 30)) {
		t.Fatalf("want only Friday fire, got %+v", fires)
	}
}

func TestOffWorkFires_EmptyRepeatMeansEveryDay(t *testing.T) {
	data := baseData()
	data.OffWork.RepeatDays = nil
	now := at(2026, time.October, 17, 10, 0) // 周六
	fires := NextFires(data, now, 60*time.Hour)
	// 周六、周日、周一都应命中（end=周一 22:00，含周一 17:30）。
	if len(fires) != 3 {
		t.Fatalf("want 3 fires (every day), got %d: %+v", len(fires), fires)
	}
}

func TestOffWorkFires_AdanceBeforeMain(t *testing.T) {
	data := baseData()
	data.OffWork.Enabled = true
	data.Remind.AdvanceMinutes = 10
	now := at(2026, time.October, 14, 10, 0)
	fires := NextFires(data, now, 24*time.Hour)
	if len(fires) != 2 {
		t.Fatalf("want 2 fires (advance + main), got %d: %+v", len(fires), fires)
	}
	if fires[0].Kind != RemindKindOffWorkAdvance || fires[1].Kind != RemindKindOffWork {
		t.Errorf("order wrong: %+v", fires)
	}
	if !fires[0].At.Equal(at(2026, time.October, 14, 17, 20)) {
		t.Errorf("advance time = %v, want 17:20", fires[0].At)
	}
}

func TestOffWorkFires_Disabled(t *testing.T) {
	data := baseData()
	data.OffWork.Enabled = false
	now := at(2026, time.October, 14, 10, 0)
	if fires := NextFires(data, now, 24*time.Hour); len(fires) != 0 {
		t.Errorf("disabled config should yield no fires, got %+v", fires)
	}
}

func TestTodoFires(t *testing.T) {
	now := at(2026, time.October, 14, 10, 0)
	data := AppData{
		OffWork: OffWorkConfig{Enabled: false},
		Todos: []Todo{
			{ID: "a", Deadline: now.Add(time.Hour).UnixMilli()},                  // 未来 → 触发
			{ID: "b", Deadline: now.Add(time.Hour).UnixMilli(), Completed: true}, // 完成 → 忽略
			{ID: "c", Deadline: now.Add(time.Hour).UnixMilli(), Reminded: true},  // 已提醒 → 忽略
			{ID: "d", Deadline: now.Add(-time.Hour).UnixMilli()},                 // 过去 → 忽略
		},
	}
	fires := NextFires(data, now, 24*time.Hour)
	if len(fires) != 1 {
		t.Fatalf("want 1 todo fire, got %d: %+v", len(fires), fires)
	}
	if fires[0].TodoID != "a" || fires[0].Kind != RemindKindTodo {
		t.Errorf("got %+v, want todo a", fires[0])
	}
}

// 代办与下班共用全局提前量：提醒时刻应为 deadline - advance。
func TestTodoFires_AdvanceBeforeDeadline(t *testing.T) {
	now := at(2026, time.October, 14, 10, 0)
	deadline := at(2026, time.October, 14, 11, 0) // 11:00 到期
	data := AppData{
		OffWork: OffWorkConfig{Enabled: false},
		Remind:  RemindSettings{AdvanceMinutes: 15},
		Todos:   []Todo{{ID: "a", Deadline: deadline.UnixMilli()}},
	}
	fires := NextFires(data, now, 24*time.Hour)
	if len(fires) != 1 {
		t.Fatalf("want 1 todo fire, got %d: %+v", len(fires), fires)
	}
	want := at(2026, time.October, 14, 10, 45) // 10:45 = 11:00 - 15min
	if !fires[0].At.Equal(want) || fires[0].Kind != RemindKindTodo {
		t.Errorf("got %+v, want todo@%v", fires[0], want)
	}
}

func TestNextFires_CombinedSorted(t *testing.T) {
	now := at(2026, time.October, 14, 10, 0)
	data := baseData()
	data.Todos = []Todo{
		{ID: "x", Deadline: at(2026, time.October, 14, 16, 0).UnixMilli()}, // 早于下班
		{ID: "y", Deadline: at(2026, time.October, 14, 18, 0).UnixMilli()}, // 晚于下班
	}
	fires := NextFires(data, now, 24*time.Hour)
	if len(fires) != 3 {
		t.Fatalf("want 3 fires, got %d: %+v", len(fires), fires)
	}
	for i := 1; i < len(fires); i++ {
		if fires[i-1].At.After(fires[i].At) {
			t.Fatalf("not sorted: %+v", fires)
		}
	}
	if fires[0].TodoID != "x" || fires[1].Kind != RemindKindOffWork || fires[2].TodoID != "y" {
		t.Errorf("unexpected order: %+v", fires)
	}
}
