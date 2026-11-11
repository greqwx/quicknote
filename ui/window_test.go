// Copyright (c) Greg and Contributors

package ui

import (
	"os"
	"path/filepath"
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/widget"

	"quicknote/config"
	"quicknote/notes"
)

func newTestWindow(t *testing.T) (*NoteWindow, *notes.Storage, string) {
	t.Helper()

	dir := t.TempDir()
	storage, err := notes.NewStorage(dir)
	if err != nil {
		t.Fatalf("storage: %v", err)
	}

	cfg := &config.Config{
		General: config.General{AutoSave: true},
		Tabs:    config.Tabs{ShowRecentTabs: true, RecentTabsCount: 5},
	}

	nw := NewNoteWindow(test.NewApp(), cfg, filepath.Join(dir, "config.ini"), storage)
	t.Cleanup(func() {
		if nw.win != nil {
			nw.win.Close()
		}
	})
	return nw, storage, dir
}

func walk(o fyne.CanvasObject, fn func(fyne.CanvasObject)) {
	fn(o)
	switch v := o.(type) {
	case *fyne.Container:
		for _, child := range v.Objects {
			walk(child, fn)
		}
	case *container.Scroll:
		if v.Content != nil {
			walk(v.Content, fn)
		}
	}
}

func findEntry(t *testing.T, nw *NoteWindow) *widget.Entry {
	t.Helper()
	var entry *widget.Entry
	walk(nw.win.Content(), func(o fyne.CanvasObject) {
		if v, ok := o.(*widget.Entry); ok {
			entry = v
		}
	})
	if entry == nil {
		t.Fatal("entry not found")
	}
	return entry
}

func findList(t *testing.T, nw *NoteWindow) *widget.List {
	t.Helper()
	var list *widget.List
	walk(nw.win.Content(), func(o fyne.CanvasObject) {
		if v, ok := o.(*widget.List); ok {
			list = v
		}
	})
	if list == nil {
		t.Fatal("recent notes list not found")
	}
	return list
}

func findButton(t *testing.T, nw *NoteWindow, text string) *widget.Button {
	t.Helper()
	var btn *widget.Button
	walk(nw.win.Content(), func(o fyne.CanvasObject) {
		if v, ok := o.(*widget.Button); ok && v.Text == text {
			btn = v
		}
	})
	if btn == nil {
		t.Fatalf("button %q not found", text)
	}
	return btn
}

func buttonExists(nw *NoteWindow, text string) bool {
	found := false
	walk(nw.win.Content(), func(o fyne.CanvasObject) {
		if v, ok := o.(*widget.Button); ok && v.Text == text {
			found = true
		}
	})
	return found
}

func TestAutosaveReplacesSaveWithOpenConfig(t *testing.T) {
	nw, _, _ := newTestWindow(t) // AutoSave is on there
	nw.Show()

	if buttonExists(nw, "Save") {
		t.Error("Save should be gone while AutoSave is on")
	}
	findButton(t, nw, "Open config")
	findButton(t, nw, "Clear")
}

func TestBothButtonsWithoutAutosave(t *testing.T) {
	nw, _, _ := newTestWindow(t)
	nw.cfg.General.AutoSave = false
	nw.Show()

	findButton(t, nw, "Save")
	findButton(t, nw, "Clear")
	findButton(t, nw, "Open config")
}

func TestNoteLabel(t *testing.T) {
	cases := map[string]string{
		`C:\notes\note_2026-11-10_21-59-45.txt`: "10.11 21:59",
		"/notes/note_2026-01-02_03-04-05.txt":   "02.01 03:04",
		"/notes/weird-name.txt":                 "weird-name",
	}
	for path, want := range cases {
		if got := noteLabel(path); got != want {
			t.Errorf("noteLabel(%q) = %q, want %q", path, got, want)
		}
	}
}

func TestRecentPanelIsOnTheLeft(t *testing.T) {
	nw, _, dir := newTestWindow(t)
	nw.cfg.Tabs.RecentTabsCount = 20
	for _, n := range []string{
		"note_2026-11-01_08-00-00.txt", "note_2026-11-02_09-00-00.txt", "note_2026-11-03_10-00-00.txt",
		"note_2026-11-04_11-00-00.txt", "note_2026-11-05_12-00-00.txt", "note_2026-11-06_13-00-00.txt",
		"note_2026-11-07_14-00-00.txt", "note_2026-11-08_15-00-00.txt",
	} {
		if err := os.WriteFile(filepath.Join(dir, n), []byte("x"), 0644); err != nil {
			t.Fatal(err)
		}
	}

	nw.Show()
	nw.win.Resize(fyne.NewSize(560, 360))

	list := findList(t, nw)
	if got := list.Length(); got != 8 {
		t.Errorf("list has %d notes, want 8", got)
	}
	if pos := fyne.CurrentApp().Driver().AbsolutePositionForObject(list); pos.X > 10 {
		t.Errorf("panel X = %v, want it at the left edge", pos.X)
	}
	if got := list.Size().Width; got != sidePanelWidth {
		t.Errorf("panel width = %v, want %v", got, sidePanelWidth)
	}
}

func findStatus(t *testing.T, nw *NoteWindow) *widget.Label {
	t.Helper()
	var label *widget.Label
	walk(nw.win.Content(), func(o fyne.CanvasObject) {
		if v, ok := o.(*widget.Label); ok && v.Text == "" {
			label = v
		}
	})
	if label == nil {
		t.Fatal("status label not found")
	}
	return label
}

func TestStatusAppearsOnlyWithText(t *testing.T) {
	nw, storage, _ := newTestWindow(t)
	if _, err := storage.Save("hello"); err != nil {
		t.Fatalf("save: %v", err)
	}

	nw.Show()
	status := findStatus(t, nw)
	if status.Visible() {
		t.Error("empty status label takes space")
	}

	findList(t, nw).Select(0)
	if !status.Visible() {
		t.Error("status should appear after opening a note")
	}

	findButton(t, nw, "Clear").OnTapped()
	if status.Visible() {
		t.Error("status should disappear after Clear")
	}
}

func TestEscapeClosesWindow(t *testing.T) {
	nw, storage, _ := newTestWindow(t)
	if _, err := storage.Save("hello"); err != nil {
		t.Fatalf("save: %v", err)
	}
	nw.Show()

	nw.win.Canvas().OnTypedKey()(&fyne.KeyEvent{Name: fyne.KeyEscape})

	if nw.win != nil {
		t.Error("window stayed open")
	}
}

func TestRecentListLoadsNote(t *testing.T) {
	nw, storage, _ := newTestWindow(t)
	path, err := storage.Save("original")
	if err != nil {
		t.Fatalf("save: %v", err)
	}

	nw.Show()
	findList(t, nw).Select(0)

	if got := findEntry(t, nw).Text; got != "original" {
		t.Errorf("entry = %q, want %q", got, "original")
	}
	if nw.currentPath != path {
		t.Errorf("currentPath = %q, want %q", nw.currentPath, path)
	}
}

func TestAutosaveUpdatesOpenedNote(t *testing.T) {
	nw, storage, dir := newTestWindow(t)
	path, err := storage.Save("original")
	if err != nil {
		t.Fatalf("save: %v", err)
	}

	nw.Show()
	findList(t, nw).Select(0)
	findEntry(t, nw).SetText("edited")
	nw.win.Close()

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read note: %v", err)
	}
	if string(data) != "edited" {
		t.Errorf("note content = %q, want %q", data, "edited")
	}

	files, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("read dir: %v", err)
	}
	if len(files) != 1 {
		t.Errorf("got %d notes, want 1 (no duplicate)", len(files))
	}
}

func TestAutosaveCreatesNoteForUnsavedText(t *testing.T) {
	nw, _, dir := newTestWindow(t)
	nw.Show()

	findEntry(t, nw).SetText("brand new")
	nw.win.Close()

	files, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("read dir: %v", err)
	}
	if len(files) != 1 {
		t.Fatalf("got %d notes, want 1", len(files))
	}
	data, err := os.ReadFile(dir + string(os.PathSeparator) + files[0].Name())
	if err != nil {
		t.Fatalf("read note: %v", err)
	}
	if string(data) != "brand new" {
		t.Errorf("note content = %q, want %q", data, "brand new")
	}
}
