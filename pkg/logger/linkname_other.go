//go:build !windows

package logger

// logLinkNameSupported reports whether the platform can keep the rotating log
// symlink (logs/app.log) without elevated privileges.
func logLinkNameSupported() bool { return true }