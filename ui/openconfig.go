// Copyright (c) Greg and Contributors

package ui

import (
	"errors"
	"os/exec"
	"runtime"
)

// openConfigFile hands the config file to the system so the user can edit it.
func openConfigFile(path string) error {
	if path == "" {
		return errors.New("config file path is unknown")
	}

	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("notepad.exe", path)
	case "darwin":
		cmd = exec.Command("open", path)
	default:
		cmd = exec.Command("xdg-open", path)
	}
	return cmd.Start()
}
