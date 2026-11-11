// Copyright (c) Greg and Contributors

package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

// appTheme keeps the dark palette but squares off the rounded widget corners.
type appTheme struct {
	fyne.Theme
}

func (t appTheme) Size(name fyne.ThemeSizeName) float32 {
	switch name {
	case theme.SizeNameInputRadius, theme.SizeNameSelectionRadius:
		return 0
	}
	return t.Theme.Size(name)
}

// ApplyTheme makes QuickNote dark and square even when the OS uses a light theme.
func ApplyTheme(a fyne.App) {
	a.Settings().SetTheme(appTheme{theme.DarkTheme()})
}
