// Copyright (c) Greg and Contributors

package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoadCreatesFullConfig(t *testing.T) {
	path := filepath.Join(t.TempDir(), "nested", "config.ini")

	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if !cfg.General.AutoSave {
		t.Errorf("AutoSave = false, want true")
	}
	if !cfg.General.Autostart {
		t.Errorf("Autostart = false, want true")
	}
	if cfg.General.NotesDirectory == "" || !strings.HasSuffix(cfg.General.NotesDirectory, "Quicknote") {
		t.Errorf("NotesDirectory = %q, want .../Quicknote", cfg.General.NotesDirectory)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read generated config: %v", err)
	}
	for _, want := range []string{"[General]", "AutoSave", "Autostart", "NotesDirectory", "[Tabs]", "ShowRecentTabs", "RecentTabsCount"} {
		if !strings.Contains(string(data), want) {
			t.Errorf("generated config is missing %q:\n%s", want, data)
		}
	}
}

func TestSaveLoadRoundTrip(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.ini")

	want := Defaults()
	want.General.AutoSave = false
	want.General.Autostart = false
	want.General.NotesDirectory = `C:\Users\Greg\Documents\Quicknote`
	want.Tabs.ShowRecentTabs = false
	want.Tabs.RecentTabsCount = 12

	if err := Save(path, want); err != nil {
		t.Fatalf("Save: %v", err)
	}

	got, err := Load(path)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if *got != *want {
		t.Errorf("round trip mismatch:\n got %+v\nwant %+v", *got, *want)
	}
}

func TestLoadKeepsExistingConfig(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.ini")
	if err := os.WriteFile(path, []byte("[General]\nAutoSave = false\nNotesDirectory = /tmp/notes\n"), 0644); err != nil {
		t.Fatal(err)
	}

	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if cfg.General.AutoSave {
		t.Errorf("AutoSave = true, want false")
	}
	if !cfg.General.Autostart {
		t.Errorf("Autostart = false, want default true")
	}
	if cfg.General.NotesDirectory != "/tmp/notes" {
		t.Errorf("NotesDirectory = %q, want /tmp/notes", cfg.General.NotesDirectory)
	}
	// Missing sections still get defaults.
	if cfg.Tabs.RecentTabsCount == 0 {
		t.Errorf("RecentTabsCount = 0, want default")
	}
}

func TestPathHonoursEnv(t *testing.T) {
	t.Setenv("QUICKNOTE_CONFIG", `C:\custom\config.ini`)

	got, err := Path()
	if err != nil {
		t.Fatalf("Path: %v", err)
	}
	if got != `C:\custom\config.ini` {
		t.Errorf("Path = %q", got)
	}
}
