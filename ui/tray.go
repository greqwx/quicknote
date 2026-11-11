// Copyright (c) Greg and Contributors

package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/desktop"

	"quicknote/assets"
	"quicknote/autostart"
	"quicknote/config"
)

// SetupTray adds the system-tray icon and menu items.
// showNote is called when the user clicks "Open Note" from the tray.
func SetupTray(app fyne.App, cfg *config.Config, cfgPath string, showNote func()) {
	desk, ok := app.(desktop.App)
	if !ok {
		// not a desktop environment, skip tray
		return
	}

	openItem := fyne.NewMenuItem("Open Note  (Ctrl+Shift+N)", func() { showNote() })

	autostartItem := fyne.NewMenuItem("Start with system", nil)
	autostartItem.Checked = autostart.Enabled()

	quitItem := fyne.NewMenuItem("Quit", func() { app.Quit() })

	menu := fyne.NewMenu("QuickNote", openItem, autostartItem, fyne.NewMenuItemSeparator(), quitItem)

	autostartItem.Action = func() {
		enable := !autostart.Enabled()
		var err error
		if enable {
			err = autostart.Enable()
		} else {
			err = autostart.Disable()
		}
		if err != nil {
			fyne.LogError("autostart", err)
		} else {
			cfg.General.Autostart = enable
			if cfgPath != "" {
				if err := config.Save(cfgPath, cfg); err != nil {
					fyne.LogError("save config", err)
				}
			}
		}
		autostartItem.Checked = autostart.Enabled()
		desk.SetSystemTrayMenu(menu)
	}

	desk.SetSystemTrayMenu(menu)
	desk.SetSystemTrayIcon(assets.TrayIcon)
}
