// Copyright (c) Greg and Contributors

package notes

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type Storage struct {
	dir string
}

func NewStorage(dir string) (*Storage, error) {
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("cannot create notes dir: %w", err)
	}
	return &Storage{dir: dir}, nil
}

// Save writes text to a timestamped file and returns its path.
func (s *Storage) Save(text string) (string, error) {
	if text == "" {
		return "", nil
	}
	name := fmt.Sprintf("note_%s.txt", time.Now().Format("2006-01-02_15-04-05"))
	path := filepath.Join(s.dir, name)
	if err := os.WriteFile(path, []byte(text), 0644); err != nil {
		return "", err
	}
	return path, nil
}

// Update overwrites an existing note with text.
func (s *Storage) Update(path, text string) error {
	return os.WriteFile(path, []byte(text), 0644)
}

// Recent returns the N most recently modified note paths.
func (s *Storage) Recent(n int) ([]string, error) {
	entries, err := os.ReadDir(s.dir)
	if err != nil {
		return nil, err
	}
	// ReadDir returns entries sorted by name; our names are timestamps, so newest is last.
	var paths []string
	for i := len(entries) - 1; i >= 0 && len(paths) < n; i-- {
		e := entries[i]
		if !e.IsDir() {
			paths = append(paths, filepath.Join(s.dir, e.Name()))
		}
	}
	return paths, nil
}
