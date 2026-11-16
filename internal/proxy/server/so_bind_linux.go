//go:build linux

package server

import "syscall"

// soBindToDevice is SO_BINDTODEVICE on Linux.
const soBindToDevice = syscall.SO_BINDTODEVICE

func setBindToDevice(fd uintptr, iface string) error {
	return syscall.SetsockoptString(int(fd), syscall.SOL_SOCKET, soBindToDevice, iface)
}
