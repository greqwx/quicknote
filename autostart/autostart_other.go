// Copyright (c) Greg and Contributors

//go:build !windows && !linux && !darwin

package autostart

import "errors"

var errUnsupported = errors.New("autostart is not supported on this platform")

func Enable() error { return errUnsupported }

func Disable() error { return nil }

func Enabled() bool { return false }
