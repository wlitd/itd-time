// Package reminder 是提醒调度与展示的编排层（Coordinator）。
// 它不含「何时提醒」的规则（在 domain.NextFires 纯函数里），只负责：
//  1. 轮询检测到期提醒并入队；
//  2. 以「并行堆叠」方式在右下角逐个错位弹出透明置顶提醒窗（互不阻塞）；
//  3. 每条提醒独立计时/响应，弹出后回写代办提醒态，避免重启后重复。
//
// 错过的提醒（应用未运行期间到期的）因轮询只看未来窗口而自然丢弃，符合既定策略。
package reminder

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
	"sync"
	"time"

	"github.com/wailsapp/wails/v3/pkg/application"

	"itdtime/internal/domain"
	"itdtime/internal/storage"
)

const (
	// pollInterval 到期检测粒度；与 V2 一致的每秒轮询，但逻辑落在可测的 domain 层。
	pollInterval = time.Second
	// autoCloseDelay 单条提醒窗无人响应的兜底自动关闭时长（各窗独立计时，互不影响）。
	autoCloseDelay = 20 * time.Second
	// popupWidth/popupHeight/margin 右下角提醒卡的尺寸与屏幕边距。
	// 卡片为单行布局：头像 + 标题/描述 + 双按钮同行（按钮文案走 Tooltip）。
	popupWidth  = 380
	popupHeight = 64
	// popupMargin 提醒窗与屏幕右缘的间距，取较小值使其更贴靠右下角。
	popupMargin = 6
	// popupBottomMargin 底部与任务栏的贴合间距，小于侧边距使提醒窗更靠下。
	popupBottomMargin = 6
	// popupGap 堆叠时相邻两张卡片之间的竖直间距。
	popupGap = 8
	// maxStack 同时在屏的提醒上限；超出时淘汰最旧一条腾出位置，避免无限堆叠。
	maxStack = 6
)

// Request 一次待展示的提醒。
type Request struct {
	Kind   domain.RemindKind
	Title  string
	TodoID string // 仅代办类提醒有效，用于回写 Reminded。
}

// entry 一条在屏提醒：唯一 id（供前端 Dismiss 精确定位）+ 窗口句柄。
type entry struct {
	id  int
	win *application.WebviewWindow
}

// layoutTarget 一次重排计算出的目标坐标（在锁外执行 SetPosition，避免持锁调用原生窗口）。
type layoutTarget struct {
	win *application.WebviewWindow
	x   int
	y   int
}

// Coordinator 提醒编排器，作为一个绑定 Service 暴露 Dismiss/Preview 给前端。
type Coordinator struct {
	store *storage.Store
	app   *application.App

	queue chan Request  // 待展示提醒入口（有界缓冲，满则丢弃最新项，绝不阻塞调度）。
	stop  chan struct{} // 停止调度与展示循环。

	mu     sync.Mutex // 保护 stack 与 nextID。
	nextID int        // 自增的提醒窗 id。
	stack  []*entry   // 当前在屏提醒，index 0 = 最底部（最早），越靠后越往上。
}

// NewCoordinator 构造 Coordinator；依赖存储（读配置/回写提醒态）。
func NewCoordinator(store *storage.Store) *Coordinator {
	return &Coordinator{
		store: store,
		queue: make(chan Request, 64),
		stop:  make(chan struct{}),
	}
}

// ServiceStartup 实现 application.ServiceStartup：拿到 App 后启动轮询与展示循环。
func (c *Coordinator) ServiceStartup(_ context.Context, _ application.ServiceOptions) error {
	c.app = application.Get()
	go c.pollLoop()
	go c.displayLoop()
	return nil
}

// ServiceShutdown 实现 application.ServiceShutdown：优雅停止后台 goroutine。
func (c *Coordinator) ServiceShutdown() error {
	close(c.stop)
	return nil
}

// pollLoop 每秒以「刚过去的时间窗」调用 domain.NextFires，捕捉落入该窗口的提醒并入队。
func (c *Coordinator) pollLoop() {
	ticker := time.NewTicker(pollInterval)
	defer ticker.Stop()
	prev := time.Now()
	for {
		select {
		case now := <-ticker.C:
			elapsed := now.Sub(prev)
			for _, fire := range domain.NextFires(c.store.Data(), prev, elapsed) {
				c.enqueue(fire)
			}
			prev = now
		case <-c.stop:
			return
		}
	}
}

// enqueue 把一次触发解析为可读文案并入队；队列满时丢弃（不阻塞调度循环）。
func (c *Coordinator) enqueue(fire domain.Fire) {
	req := Request{Kind: fire.Kind, TodoID: fire.TodoID}
	switch fire.Kind {
	case domain.RemindKindOffWork:
		req.Title = "下班啦，辛苦了！"
	case domain.RemindKindOffWorkAdvance:
		req.Title = "快到下班时间了"
	case domain.RemindKindTodo:
		req.Title = c.todoTitle(fire.TodoID)
	}
	select {
	case c.queue <- req:
	default:
	}
}

// displayLoop 取出提醒即时弹出：并行堆叠，处理下一条不等待上一条关闭。
func (c *Coordinator) displayLoop() {
	for {
		select {
		case req := <-c.queue:
			c.showOne(req)
		case <-c.stop:
			return
		}
	}
}

// showOne 即时展示一条提醒，并为其单独安排兜底自动关闭。
func (c *Coordinator) showOne(req Request) {
	// 代办提醒：展示即标记已提醒，确保「重启不重复」。
	if req.TodoID != "" {
		c.markReminded(req.TodoID)
	}
	e := c.openPopup(req)
	if e == nil {
		return
	}
	// 每条独立兜底：到 autoCloseDelay 自动关闭；若在此之前已被 Dismiss 则 closeEntry 幂等跳过。
	go func(id int) {
		timer := time.NewTimer(autoCloseDelay)
		defer timer.Stop()
		select {
		case <-timer.C:
			c.closeEntry(id)
		case <-c.stop:
		}
	}(e.id)
}

// openPopup 创建一个右下角透明置顶提醒窗并堆叠到栈顶；返回其在屏 entry。
// 内容（含唯一 id）通过 URL query 传递，供提醒窗回调 Dismiss 精确定位。
func (c *Coordinator) openPopup(req Request) *entry {
	if c.app == nil {
		return nil
	}

	// 分配 id 并在超出上限时摘除最旧一条（仅从 stack 移除，关窗动作放到锁外）。
	c.mu.Lock()
	c.nextID++
	id := c.nextID
	var evicted *application.WebviewWindow
	if len(c.stack) >= maxStack {
		oldest := c.stack[0]
		c.stack = c.stack[1:]
		evicted = oldest.win
	}
	c.mu.Unlock()

	q := url.Values{}
	q.Set("remind", "1")
	q.Set("kind", string(req.Kind))
	q.Set("title", req.Title)
	q.Set("id", strconv.Itoa(id))

	win := c.app.Window.NewWithOptions(application.WebviewWindowOptions{
		Name:           fmt.Sprintf("remind-%d", id),
		Width:          popupWidth,
		Height:         popupHeight,
		Frameless:      true,
		AlwaysOnTop:    true,
		DisableResize:  true,
		BackgroundType: application.BackgroundTypeTransparent,
		Hidden:         true,
		URL:            "/?" + q.Encode(),
	})
	e := &entry{id: id, win: win}

	// 入栈并计算整叠窗口的目标坐标（含新窗与各旧窗的上移补位）。
	c.mu.Lock()
	c.stack = append(c.stack, e)
	targets := c.layoutLocked()
	c.mu.Unlock()

	// 锁外执行原生窗口操作：先就位坐标，再显示新窗，最后关闭被淘汰的最旧窗。
	for _, t := range targets {
		t.win.SetPosition(t.x, t.y)
	}
	win.Show()
	if evicted != nil {
		evicted.Close()
	}
	return e
}

// Dismiss 供提醒窗的关闭按钮调用（绑定方法）：按 id 关闭对应提醒并重排其余补位。
func (c *Coordinator) Dismiss(id int) {
	c.closeEntry(id)
}

// Preview 供设置页「预览提醒样式」调用（绑定方法）：立即按给定样式排入一条提醒。
func (c *Coordinator) Preview(kind domain.RemindKind, title string) {
	req := Request{Kind: kind, Title: title}
	if req.Title == "" {
		req.Title = "这是一条提醒预览"
	}
	select {
	case c.queue <- req:
	default:
	}
}

// closeEntry 关闭指定 id 的提醒并从栈中移除，剩余各窗重排补位；对未知/已关闭 id 幂等。
func (c *Coordinator) closeEntry(id int) {
	c.mu.Lock()
	idx := -1
	var win *application.WebviewWindow
	for i, e := range c.stack {
		if e.id == id {
			idx = i
			win = e.win
			break
		}
	}
	if idx < 0 {
		c.mu.Unlock()
		return
	}
	c.stack = append(c.stack[:idx], c.stack[idx+1:]...)
	targets := c.layoutLocked()
	c.mu.Unlock()

	win.Close()
	for _, t := range targets {
		t.win.SetPosition(t.x, t.y)
	}
}

// layoutLocked 依据当前 stack 顺序计算各提醒窗自底向上的堆叠坐标。调用方须持有 c.mu。
func (c *Coordinator) layoutLocked() []layoutTarget {
	screen := c.app.Screen.GetPrimary()
	if screen == nil {
		return nil
	}
	area := screen.WorkArea
	baseX := area.X + area.Width - popupWidth - popupMargin
	if baseX < area.X {
		baseX = area.X
	}
	bottomY := area.Y + area.Height - popupHeight - popupBottomMargin
	if bottomY < area.Y {
		bottomY = area.Y
	}
	targets := make([]layoutTarget, 0, len(c.stack))
	for i, e := range c.stack {
		y := bottomY - i*(popupHeight+popupGap)
		if y < area.Y {
			y = area.Y
		}
		targets = append(targets, layoutTarget{win: e.win, x: baseX, y: y})
	}
	return targets
}

// todoTitle 按 ID 查代办标题；查不到时给兜底文案。
func (c *Coordinator) todoTitle(id string) string {
	for _, t := range c.store.Data().Todos {
		if t.ID == id {
			return t.Title
		}
	}
	return "有一条代办到期了"
}

// markReminded 回写代办提醒态，落盘持久化。
func (c *Coordinator) markReminded(id string) {
	_ = c.store.Update(func(data *domain.AppData) error {
		for i := range data.Todos {
			if data.Todos[i].ID == id {
				data.Todos[i].Reminded = true
			}
		}
		return nil
	})
}
