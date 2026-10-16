// Package storage 负责将 domain.AppData 持久化到本地单一 JSON 文件。
// 写入采用「临时文件 + 原子 rename」，保证崩溃/断电后不产生半截文件。
package storage

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"itdtime/internal/domain"
)

// Store 是应用数据的持久化边界，进程内持有一份内存副本。
// 非并发安全：所有写操作都应经由上层的 Service 串行化调用。
type Store struct {
	path string
	data domain.AppData
}

// New 基于用户配置目录构建 Store。目录布局：
//   - Windows: %APPDATA%\itdtime\data.json
//   - macOS:   ~/Library/Application Support/itdtime/data.json
//   - Linux:   ~/.config/itdtime/data.json
func New() (*Store, error) {
	base, err := os.UserConfigDir()
	if err != nil {
		return nil, fmt.Errorf("resolve user config dir: %w", err)
	}
	dir := filepath.Join(base, "itdtime")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, fmt.Errorf("create config dir: %w", err)
	}
	return &Store{path: filepath.Join(dir, "data.json")}, nil
}

// Load 读取磁盘数据到内存。文件不存在时载入默认值（不算错误）。
func (s *Store) Load() error {
	raw, err := os.ReadFile(s.path)
	if errors.Is(err, os.ErrNotExist) {
		s.data = domain.DefaultAppData()
		return s.flush()
	}
	if err != nil {
		return fmt.Errorf("read data file: %w", err)
	}

	// 以默认值为基底再覆盖，兼容旧版本缺失字段。
	merged := domain.DefaultAppData()
	if err := json.Unmarshal(raw, &merged); err != nil {
		return fmt.Errorf("parse data json: %w", err)
	}
	if merged.Todos == nil {
		merged.Todos = []domain.Todo{}
	}
	if err := s.migrate(&merged); err != nil {
		return err
	}
	s.data = merged
	return nil
}

// migrate 在数据结构版本演进时做迁移。当前仅基线 v1。
func (s *Store) migrate(data *domain.AppData) error {
	if data.SchemaVersion == 0 {
		data.SchemaVersion = 1
	}
	// 未来：switch data.SchemaVersion { case 1: ... }
	return nil
}

// Data 返回内存副本快照（值拷贝，todos 单独切片拷贝以避免外部改写）。
func (s *Store) Data() domain.AppData {
	out := s.data
	out.Todos = append([]domain.Todo(nil), s.data.Todos...)
	return out
}

// Update 以回调方式修改数据：在副本上试改，仅当回调成功（返回 nil）时
// 才提交到内存并原子落盘；回调返回错误则内存与磁盘均保持不变。
func (s *Store) Update(fn func(data *domain.AppData) error) error {
	candidate := s.data
	candidate.Todos = append([]domain.Todo(nil), s.data.Todos...)
	if err := fn(&candidate); err != nil {
		return err
	}
	prev := s.data
	s.data = candidate
	if err := s.flush(); err != nil {
		s.data = prev // 落盘失败回滚内存，保持一致性。
		return err
	}
	return nil
}

// flush 将当前内存数据原子写入磁盘。
func (s *Store) flush() error {
	raw, err := json.MarshalIndent(s.data, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal data: %w", err)
	}
	tmp := s.path + ".tmp"
	if err := os.WriteFile(tmp, raw, 0o644); err != nil {
		return fmt.Errorf("write temp file: %w", err)
	}
	if err := os.Rename(tmp, s.path); err != nil {
		return fmt.Errorf("atomic rename: %w", err)
	}
	return nil
}
