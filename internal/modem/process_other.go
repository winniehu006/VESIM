//go:build !windows

package modem

import "syscall"

// terminateProcess asks the process to exit. Used by the Linux fuser-based port
// cleanup path.
func terminateProcess(pid int) error {
	return syscall.Kill(pid, syscall.SIGTERM)
}
