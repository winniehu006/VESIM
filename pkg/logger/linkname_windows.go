//go:build windows

package logger

// logLinkNameSupported reports whether the platform can keep the rotating log
// symlink (logs/app.log) without elevated privileges. Windows requires
// SeCreateSymbolicLinkPrivilege, which normal users do not have.
func logLinkNameSupported() bool { return false }