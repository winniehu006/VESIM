//go:build !linux

package netprobe

// setBindToDevice is a no-op outside Linux (SO_BINDTODEVICE is Linux-only).
func setBindToDevice(_ uintptr, _ string) error {
	return nil
}
