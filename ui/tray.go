// Copyright (c) Greg and Contributors

package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/desktop"

	"quicknote/assets"
)

// SetupTray adds the system-tray icon and menu items.
// showNote is called when the user clicks "Open Note" from the tray.
func SetupTray(app fyne.App, showNote func()) {
	desk, ok := app.(desktop.App)
	if !ok {
		// not a desktop environment, skip tray
		return
	}

	openItem := fyne.NewMenuItem("Open Note  (Ctrl+Shift+N)", func() { showNote() })
	quitItem := fyne.NewMenuItem("Quit", func() { app.Quit() })

	menu := fyne.NewMenu("QuickNote", openItem, fyne.NewMenuItemSeparator(), quitItem)
	desk.SetSystemTrayMenu(menu)
	desk.SetSystemTrayIcon(assets.TrayIcon)
}
