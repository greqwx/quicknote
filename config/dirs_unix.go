// Copyright (c) Greg and Contributors

//go:build !windows

package config

import (
	"os"
	"path/filepath"
	"strings"
)

func documentsDir() string {
	home, _ := os.UserHomeDir()
	if p := xdgDocumentsDir(home); p != "" {
		return p
	}
	if home == "" {
		return "Documents"
	}
	return filepath.Join(home, "Documents")
}

// xdgDocumentsDir honours a relocated Documents folder on Linux desktops.
// macOS has no user-dirs.dirs, so this simply returns "".
func xdgDocumentsDir(home string) string {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return ""
	}
	data, err := os.ReadFile(filepath.Join(configDir, "user-dirs.dirs"))
	if err != nil {
		return ""
	}
	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if !strings.HasPrefix(line, "XDG_DOCUMENTS_DIR=") {
			continue
		}
		dir := strings.Trim(strings.TrimPrefix(line, "XDG_DOCUMENTS_DIR="), `"`)
		dir = strings.ReplaceAll(dir, "$HOME", home)
		if dir != "" && dir != home {
			return dir
		}
	}
	return ""
}
