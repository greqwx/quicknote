// Copyright (c) Greg and Contributors

//go:build windows

package config

import (
	"os"
	"path/filepath"

	"golang.org/x/sys/windows"
)

// documentsDir asks the shell for the Documents folder, which Windows may
// redirect (OneDrive, localised names). Falls back to %USERPROFILE%\Documents.
func documentsDir() string {
	if p, err := windows.KnownFolderPath(windows.FOLDERID_Documents, windows.KF_FLAG_DEFAULT); err == nil && p != "" {
		return p
	}
	if home, err := os.UserHomeDir(); err == nil {
		return filepath.Join(home, "Documents")
	}
	return "Documents"
}
