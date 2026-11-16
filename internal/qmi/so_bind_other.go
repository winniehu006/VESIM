//go:build !linux

package qmicore

// setBindToDevice is a no-op outside Linux (SO_BINDTODEVICE is Linux-only).
func setBindToDevice(_ uintptr, _ string) error {
	return nil
}
