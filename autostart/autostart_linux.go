// Copyright (c) Greg and Contributors

//go:build linux

package autostart

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const desktopEntry = `[Desktop Entry]
Type=Application
Name=Quicknote
Exec=%s
Terminal=false
X-GNOME-Autostart-enabled=true
`

// Enable drops a .desktop file into the XDG autostart directory.
func Enable() error {
	exe, err := os.Executable()
	if err != nil {
		return err
	}

	path, err := entryPath()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}

	body := fmt.Sprintf(desktopEntry, execQuoted(exe))
	return os.WriteFile(path, []byte(body), 0644)
}

// Disable deletes the autostart file; it does nothing when the file is absent.
func Disable() error {
	path, err := entryPath()
	if err != nil {
		return err
	}
	if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}

// Enabled reports whether the autostart file exists.
func Enabled() bool {
	path, err := entryPath()
	if err != nil {
		return false
	}
	_, err = os.Stat(path)
	return err == nil
}

func entryPath() (string, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "autostart", "quicknote.desktop"), nil
}

// execQuoted escapes a path for the Exec key, which is not run through a shell
// but still has its own quoting rules.
func execQuoted(path string) string {
	escaped := strings.NewReplacer(`\`, `\\`, `"`, `\"`, "`", "\\`", "$", `\$`).Replace(path)
	return `"` + escaped + `"`
}
