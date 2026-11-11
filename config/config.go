// Copyright (c) Greg and Contributors

package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/go-ini/ini"
)

const (
	appVendor = "Greg"
	appName   = "Quicknote"
	fileName  = "config.ini"
)

type General struct {
	AutoSave       bool
	NotesDirectory string
}

type Tabs struct {
	ShowRecentTabs  bool
	RecentTabsCount int
}

type Config struct {
	General General
	Tabs    Tabs
}

// Defaults is what a fresh install starts with.
func Defaults() *Config {
	return &Config{
		General: General{
			AutoSave:       true,
			NotesDirectory: filepath.Join(documentsDir(), appName),
		},
		Tabs: Tabs{
			ShowRecentTabs:  true,
			RecentTabsCount: 5,
		},
	}
}

// Path returns the config location: $QUICKNOTE_CONFIG when set, otherwise the
// per-user config directory (%AppData%\Greg\Quicknote on Windows,
// ~/.config/Greg/Quicknote on Linux, ~/Library/Application Support on macOS).
func Path() (string, error) {
	if p := os.Getenv("QUICKNOTE_CONFIG"); p != "" {
		return p, nil
	}

	dir, err := os.UserConfigDir()
	if err != nil {
		// Portable fallback: next to the executable.
		exe, exeErr := os.Executable()
		if exeErr != nil {
			return fileName, nil
		}
		return filepath.Join(filepath.Dir(exe), fileName), nil
	}
	return filepath.Join(dir, appVendor, appName, fileName), nil
}

// Load reads the config, creating a complete default file when it is missing.
// A returned config is never nil when error is nil; on error callers may fall
// back to Defaults().
func Load(path string) (*Config, error) {
	if _, err := os.Stat(path); err != nil {
		if !os.IsNotExist(err) {
			return nil, err
		}
		cfg := Defaults()
		if err := Save(path, cfg); err != nil {
			return cfg, fmt.Errorf("create default config %s: %w", path, err)
		}
		return cfg, nil
	}

	f, err := ini.Load(path)
	if err != nil {
		return nil, err
	}

	cfg := Defaults()
	if sec, err := f.GetSection("General"); err == nil {
		if k, err := sec.GetKey("AutoSave"); err == nil {
			cfg.General.AutoSave, _ = k.Bool()
		}
		if k, err := sec.GetKey("NotesDirectory"); err == nil {
			if v := k.String(); v != "" {
				cfg.General.NotesDirectory = v
			}
		}
	}

	if sec, err := f.GetSection("Tabs"); err == nil {
		if k, err := sec.GetKey("ShowRecentTabs"); err == nil {
			cfg.Tabs.ShowRecentTabs, _ = k.Bool()
		}
		if k, err := sec.GetKey("RecentTabsCount"); err == nil {
			cfg.Tabs.RecentTabsCount, _ = k.Int()
		}
	}

	return cfg, nil
}

// Save writes every field, so the file always contains the full set of keys.
func Save(path string, cfg *Config) error {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}

	f := ini.Empty()

	general := f.Section("General")
	general.NewKey("AutoSave", strconv.FormatBool(cfg.General.AutoSave)) //nolint:errcheck
	general.NewKey("NotesDirectory", cfg.General.NotesDirectory)         //nolint:errcheck

	tabs := f.Section("Tabs")
	tabs.NewKey("ShowRecentTabs", strconv.FormatBool(cfg.Tabs.ShowRecentTabs)) //nolint:errcheck
	tabs.NewKey("RecentTabsCount", strconv.Itoa(cfg.Tabs.RecentTabsCount))     //nolint:errcheck

	return f.SaveTo(path)
}
