// Copyright (c) Greg and Contributors

package assets

import (
	_ "embed"

	"fyne.io/fyne/v2"
)

//go:embed icon-128.png
var appIcon []byte

//go:embed icon.png
var trayIcon []byte

// AppIcon is the window and taskbar icon.
var AppIcon = fyne.NewStaticResource("icon.png", appIcon)

// TrayIcon is the notification-area icon. Windows needs a raster image here,
// Fyne converts it to an icon of the size the image provides.
var TrayIcon = fyne.NewStaticResource("tray.png", trayIcon)
