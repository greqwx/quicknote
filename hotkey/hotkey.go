// Copyright (c) Greg and Contributors

package hotkey

import (
	"log"

	"golang.design/x/hotkey"
)

// Handler is called when the hotkey fires.
type Handler func()

// Register binds Ctrl+Shift+N globally and calls h on each press.
// Returns a cleanup func that unregisters the hotkey.
func Register(h Handler) (func(), error) {
	hk := hotkey.New([]hotkey.Modifier{hotkey.ModCtrl, hotkey.ModShift}, hotkey.KeyN)
	if err := hk.Register(); err != nil {
		return nil, err
	}

	go func() {
		for range hk.Keydown() {
			h()
		}
	}()

	stop := func() {
		if err := hk.Unregister(); err != nil {
			log.Printf("hotkey unregister: %v", err)
		}
	}
	return stop, nil
}
