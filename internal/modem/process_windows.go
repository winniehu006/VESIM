//go:build windows

package modem

import "os"

// terminateProcess force-terminates the process on Windows, where SIGTERM is not
// available.
func terminateProcess(pid int) error {
	proc, err := os.FindProcess(pid)
	if err != nil {
		return err
	}
	return proc.Kill()
}
