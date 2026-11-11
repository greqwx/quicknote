// Copyright (c) Greg and Contributors

package main

import (
	"log"

	"fyne.io/fyne/v2/app"

	"quicknote/assets"
	"quicknote/config"
	"quicknote/hotkey"
	"quicknote/notes"
	"quicknote/ui"
)

func main() {
	cfgPath, err := config.Path()
	if err != nil {
		log.Printf("config path: %v", err)
	}
	cfg, err := config.Load(cfgPath)
	if err != nil {
		log.Printf("config: %v (using defaults)", err)
		cfg = config.Defaults()
	}

	storage, err := notes.NewStorage(cfg.General.NotesDirectory)
	if err != nil {
		log.Fatalf("storage: %v", err)
	}

	a := app.NewWithID("io.github.quicknote")
	ui.ApplyTheme(a)
	a.SetIcon(assets.AppIcon)

	noteWin := ui.NewNoteWindow(a, cfg, storage)
	showNote := func() { noteWin.Show() }

	ui.SetupTray(a, showNote)

	stopHotkey, err := hotkey.Register(showNote)
	if err != nil {
		log.Printf("hotkey unavailable: %v (use tray icon instead)", err)
	} else {
		defer stopHotkey()
	}

	// keep the app alive in the background; Fyne exits only when Quit is called
	a.Run()
}
