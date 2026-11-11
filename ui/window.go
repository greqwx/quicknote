// Copyright (c) Greg and Contributors

package ui

import (
	"os"
	"path/filepath"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"quicknote/config"
	"quicknote/notes"
)

const sidePanelWidth float32 = 170

// NoteWindow is the quick-entry popup.
type NoteWindow struct {
	app         fyne.App
	cfg         *config.Config
	storage     *notes.Storage
	win         fyne.Window
	currentPath string
}

func NewNoteWindow(app fyne.App, cfg *config.Config, storage *notes.Storage) *NoteWindow {
	return &NoteWindow{app: app, cfg: cfg, storage: storage}
}

// Show brings the window to the front, creating it if needed.
func (nw *NoteWindow) Show() {
	if nw.win != nil {
		nw.win.Show()
		nw.win.RequestFocus()
		return
	}
	nw.build()
	nw.win.Show()
}

func (nw *NoteWindow) build() {
	nw.win = nw.app.NewWindow("Quick Note")
	nw.win.Resize(fyne.NewSize(560, 360))
	nw.win.SetFixedSize(false)
	nw.win.CenterOnScreen()

	entry := widget.NewMultiLineEntry()
	entry.SetPlaceHolder("Type your note here...")
	entry.Wrapping = fyne.TextWrapWord

	status := widget.NewLabel("")
	status.Hide()
	setStatus := func(text string) {
		status.SetText(text)
		if text == "" {
			status.Hide()
			return
		}
		status.Show()
	}

	var refreshRecent func()
	save := func() {
		path, err := nw.commit(entry.Text)
		if err != nil {
			setStatus("Error: " + err.Error())
			return
		}
		if path != "" {
			setStatus("Saved: " + filepath.Base(path))
		}
		nw.currentPath = ""
		entry.SetText("")
		if refreshRecent != nil {
			refreshRecent()
		}
	}

	saveBtn := widget.NewButton("Save", save)

	clearBtn := widget.NewButton("Clear", func() {
		entry.SetText("")
		setStatus("")
		nw.currentPath = ""
	})

	buttons := container.NewGridWithColumns(2, saveBtn, clearBtn)
	bottom := container.NewVBox(buttons, status)

	editor := container.NewBorder(nil, bottom, nil, nil, entry)
	content := fyne.CanvasObject(editor)
	if nw.cfg.Tabs.ShowRecentTabs {
		if side, refresh := nw.recentPanel(entry, setStatus); side != nil {
			content = container.NewBorder(nil, nil, side, nil, editor)
			refreshRecent = refresh
		}
	}

	nw.win.SetContent(content)

	nw.win.Canvas().SetOnTypedKey(func(e *fyne.KeyEvent) {
		if e.Name == fyne.KeyEscape {
			nw.win.Close()
			return
		}
		// Ctrl+S saves without closing
		if e.Name == fyne.KeyS && nw.win.Canvas().Focused() == nil {
			save()
		}
	})

	nw.win.SetOnClosed(func() {
		if nw.cfg.General.AutoSave && entry.Text != "" {
			nw.commit(entry.Text) //nolint:errcheck
		}
		nw.currentPath = ""
		nw.win = nil
	})
}

// commit writes text to the note that is open, or to a new note when none is.
func (nw *NoteWindow) commit(text string) (string, error) {
	if text == "" {
		return "", nil
	}
	if nw.currentPath != "" {
		return nw.currentPath, nw.storage.Update(nw.currentPath, text)
	}
	return nw.storage.Save(text)
}

// recentPanel lists the newest notes down the left side, like a chat list.
func (nw *NoteWindow) recentPanel(entry *widget.Entry, setStatus func(string)) (fyne.CanvasObject, func()) {
	recent, err := nw.storage.Recent(nw.cfg.Tabs.RecentTabsCount)
	if err != nil || len(recent) == 0 {
		return nil, nil
	}

	labels := recentLabels(recent)

	list := widget.NewList(
		func() int { return len(recent) },
		func() fyne.CanvasObject { return widget.NewLabel("") },
		func(id widget.ListItemID, o fyne.CanvasObject) {
			o.(*widget.Label).SetText(labels[id])
		},
	)
	list.OnSelected = func(id widget.ListItemID) {
		data, err := os.ReadFile(recent[id])
		if err != nil {
			setStatus("Error: " + err.Error())
			return
		}
		entry.SetText(string(data))
		nw.currentPath = recent[id]
		setStatus("Editing: " + labels[id])
	}

	refresh := func() {
		paths, err := nw.storage.Recent(nw.cfg.Tabs.RecentTabsCount)
		if err != nil {
			setStatus("Error: " + err.Error())
			return
		}
		recent = paths
		labels = recentLabels(recent)
		list.UnselectAll()
		list.Refresh()
	}

	return container.New(sideLayout{width: sidePanelWidth}, list), refresh
}

func recentLabels(paths []string) []string {
	labels := make([]string, len(paths))
	for i, p := range paths {
		labels[i] = noteLabel(p)
	}
	return labels
}

// sideLayout gives the recent-notes list a fixed column on the left.
type sideLayout struct {
	width float32
}

func (l sideLayout) MinSize(objects []fyne.CanvasObject) fyne.Size {
	min := fyne.NewSize(l.width, 0)
	for _, o := range objects {
		child := o.MinSize()
		if child.Width > min.Width {
			min.Width = child.Width
		}
		if child.Height > min.Height {
			min.Height = child.Height
		}
	}
	min.Width = l.width
	return min
}

func (l sideLayout) Layout(objects []fyne.CanvasObject, size fyne.Size) {
	for _, o := range objects {
		o.Move(fyne.NewPos(0, 0))
		o.Resize(size)
	}
}

// noteLabel turns note_2026-11-10_21-59-45.txt into "10.11 21:59".
func noteLabel(path string) string {
	name := strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
	name = strings.TrimPrefix(name, "note_")
	if ts, err := time.Parse("2006-01-02_15-04-05", name); err == nil {
		return ts.Format("02.01 15:04")
	}
	return name
}
