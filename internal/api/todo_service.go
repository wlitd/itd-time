package api

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"itdtime/internal/domain"
	"itdtime/internal/storage"
)

// TodoInput 新增/修改代办时的可写字段集合（ID、完成态、提醒态由后端管理）。
type TodoInput struct {
	Title    string `json:"title"`
	Deadline int64  `json:"deadline"`
}

// TodoService 暴露代办 CRUD。
type TodoService struct {
	store *storage.Store
}

// NewTodoService 构造 TodoService。
func NewTodoService(store *storage.Store) *TodoService {
	return &TodoService{store: store}
}

// List 返回全部代办。
func (s *TodoService) List() []domain.Todo {
	return s.store.Data().Todos
}

// Add 新增一条代办。
func (s *TodoService) Add(input TodoInput) (domain.Todo, error) {
	if input.Title == "" {
		return domain.Todo{}, fmt.Errorf("todo title is required")
	}
	item := domain.Todo{
		ID:        newID(),
		Title:     input.Title,
		Deadline:  input.Deadline,
		Completed: false,
		Reminded:  false,
	}
	err := s.store.Update(func(data *domain.AppData) error {
		data.Todos = append(data.Todos, item)
		return nil
	})
	if err != nil {
		return domain.Todo{}, err
	}
	// 代办变更影响下次提醒时间，reminder.Coordinator 接入后在此触发重排（M4）。
	return item, nil
}

// Update 按 ID 覆盖标题与截止时间；截止时间变化时重置提醒态。
func (s *TodoService) Update(id string, input TodoInput) error {
	return s.store.Update(func(data *domain.AppData) error {
		for i := range data.Todos {
			if data.Todos[i].ID == id {
				if data.Todos[i].Deadline != input.Deadline {
					data.Todos[i].Reminded = false
				}
				data.Todos[i].Title = input.Title
				data.Todos[i].Deadline = input.Deadline
				return nil
			}
		}
		return fmt.Errorf("todo %q not found", id)
	})
}

// SetCompleted 切换完成态。
func (s *TodoService) SetCompleted(id string, completed bool) error {
	return s.store.Update(func(data *domain.AppData) error {
		for i := range data.Todos {
			if data.Todos[i].ID == id {
				data.Todos[i].Completed = completed
				return nil
			}
		}
		return fmt.Errorf("todo %q not found", id)
	})
}

// Delete 按 ID 移除。
func (s *TodoService) Delete(id string) error {
	return s.store.Update(func(data *domain.AppData) error {
		filtered := data.Todos[:0]
		for _, t := range data.Todos {
			if t.ID != id {
				filtered = append(filtered, t)
			}
		}
		data.Todos = filtered
		return nil
	})
}

// newID 生成 16 位十六进制随机 ID（含时间前缀以缓解极端碰撞）。
func newID() string {
	b := make([]byte, 8)
	_, _ = rand.Read(b)
	return fmt.Sprintf("%x%s", time.Now().UnixMilli()&0xffffff, hex.EncodeToString(b))
}
