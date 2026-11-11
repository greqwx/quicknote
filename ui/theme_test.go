// Copyright (c) Greg and Contributors

package ui

import (
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

func TestThemeSquaresCorners(t *testing.T) {
	th := appTheme{theme.DarkTheme()}
	for _, name := range []fyne.ThemeSizeName{theme.SizeNameInputRadius, theme.SizeNameSelectionRadius} {
		if got := th.Size(name); got != 0 {
			t.Errorf("Size(%v) = %v, want 0", name, got)
		}
	}
	if th.Size(theme.SizeNamePadding) <= 0 {
		t.Error("other sizes should come from the dark theme")
	}
}
