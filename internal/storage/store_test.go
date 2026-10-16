package storage

import (
	"os"
	"path/filepath"
	"testing"

	"itdtime/internal/domain"
)

// newTestStore 基于临时目录构造 Store，避免写入真实用户配置目录。
func newTestStore(t *testing.T) (*Store, string) {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "data.json")
	return &Store{path: path}, path
}

func TestLoad_MissingFile_WritesDefaults(t *testing.T) {
	s, path := newTestStore(t)
	if err := s.Load(); err != nil {
		t.Fatalf("Load: %v", err)
	}
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("expected data file to be created: %v", err)
	}
	data := s.Data()
	if data.SchemaVersion != 1 || data.OffWork.OffTime != "17:30" {
		t.Errorf("unexpected defaults: %+v", data)
	}
}

func TestUpdate_PersistsAndReloads(t *testing.T) {
	s, _ := newTestStore(t)
	if err := s.Load(); err != nil {
		t.Fatalf("Load: %v", err)
	}
	err := s.Update(func(d *domain.AppData) error {
		d.OffWork.OffTime = "18:00"
		d.Todos = append(d.Todos, domain.Todo{ID: "1", Title: "x", Deadline: 123})
		return nil
	})
	if err != nil {
		t.Fatalf("Update: %v", err)
	}

	// 用同一路径新建 Store 再读，验证真正落盘。
	s2 := &Store{path: s.path}
	if err := s2.Load(); err != nil {
		t.Fatalf("reload Load: %v", err)
	}
	d := s2.Data()
	if d.OffWork.OffTime != "18:00" {
		t.Errorf("offWork not persisted: %q", d.OffWork.OffTime)
	}
	if len(d.Todos) != 1 || d.Todos[0].Title != "x" {
		t.Errorf("todos not persisted: %+v", d.Todos)
	}
}

func TestUpdate_ErrorDoesNotPersist(t *testing.T) {
	s, _ := newTestStore(t)
	if err := s.Load(); err != nil {
		t.Fatalf("Load: %v", err)
	}
	before := s.Data().OffWork.OffTime
	err := s.Update(func(d *domain.AppData) error {
		d.OffWork.OffTime = "99:99"
		return os.ErrInvalid
	})
	if err == nil {
		t.Fatal("expected error from Update")
	}
	// 回调报错时内存与磁盘都应保持不变。
	if s.Data().OffWork.OffTime != before {
		t.Errorf("failed Update mutated memory; got %q want %q", s.Data().OffWork.OffTime, before)
	}
	// 磁盘仍是旧值。
	s2 := &Store{path: s.path}
	if err := s2.Load(); err != nil {
		t.Fatalf("reload: %v", err)
	}
	if s2.Data().OffWork.OffTime != before {
		t.Errorf("failed Update should not persist; got %q want %q", s2.Data().OffWork.OffTime, before)
	}
}

func TestFlush_LeavesNoTempFile(t *testing.T) {
	s, path := newTestStore(t)
	if err := s.Load(); err != nil {
		t.Fatalf("Load: %v", err)
	}
	if err := s.Update(func(d *domain.AppData) error {
		d.Remind.AdvanceMinutes = 5
		return nil
	}); err != nil {
		t.Fatalf("Update: %v", err)
	}
	if _, err := os.Stat(path + ".tmp"); err == nil {
		t.Errorf("temp file should have been renamed away")
	}
}

func TestLoad_CorruptJSON_Errors(t *testing.T) {
	s, path := newTestStore(t)
	if err := os.WriteFile(path, []byte("{not json"), 0o644); err != nil {
		t.Fatalf("write corrupt: %v", err)
	}
	if err := s.Load(); err == nil {
		t.Fatal("expected parse error on corrupt file")
	}
}

func TestData_ReturnsDefensiveCopy(t *testing.T) {
	s, _ := newTestStore(t)
	if err := s.Load(); err != nil {
		t.Fatalf("Load: %v", err)
	}
	if err := s.Update(func(d *domain.AppData) error {
		d.Todos = append(d.Todos, domain.Todo{ID: "1", Title: "keep"})
		return nil
	}); err != nil {
		t.Fatalf("Update: %v", err)
	}
	snap := s.Data()
	snap.Todos[0].Title = "mutated" // 修改快照不应影响 Store 内部数据。
	if s.Data().Todos[0].Title != "keep" {
		t.Errorf("Data() leaked internal slice; expected defensive copy")
	}
}
