//go:build windows

package ipsec

import (
	"errors"
	"net"

	"golang.org/x/sys/windows"
)

// setUDPEncap is not supported on Windows (UDP_ENCAP is Linux-only).
func setUDPEncap(conn *net.UDPConn, enable bool) error {
	return errors.New("UDP encapsulation is not supported on this platform")
}

// startErrorListener is a no-op on Windows (IP_RECVERR/MSG_ERRQUEUE are
// Linux-specific).
func (r *SocketManager) startErrorListener() {
	defer r.wg.Done()
}

// SockExtendedErr mirrors syscall.SockExtendedErr for API compatibility.
type SockExtendedErr struct {
	Errno  uint32
	Origin uint8
	Type   uint8
	Code   uint8
	Info   uint32
	Data   uint32
}

// ParseSockExtError is not supported on Windows.
func ParseSockExtError(b []byte) (*SockExtendedErr, error) {
	return nil, errors.New("extended socket errors are not supported on this platform")
}

// setSockBindToDevice is not supported on Windows (SO_BINDTODEVICE is Linux-only).
func setSockBindToDevice(fd int, device string) error {
	return errors.New("SO_BINDTODEVICE is not supported on this platform")
}

// applyReuseSocketOptions enables address reuse; SO_REUSEPORT does not exist on
// Windows.
func applyReuseSocketOptions(fd uintptr) error {
	return windows.SetsockoptInt(windows.Handle(fd), windows.SOL_SOCKET, windows.SO_REUSEADDR, 1)
}