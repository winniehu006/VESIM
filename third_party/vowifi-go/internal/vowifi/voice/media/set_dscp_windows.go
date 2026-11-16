//go:build windows

package media

import "golang.org/x/sys/windows"

const (
	ipProtoIPV6 = 41
	ipv6TClass  = 39
)

// setTrafficClass applies the IPv4 TOS value and the IPv6 traffic class value.
func setTrafficClass(fd uintptr, value int) (ipv4Err, ipv6Err error) {
	ipv4Err = windows.SetsockoptInt(windows.Handle(fd), windows.IPPROTO_IP, windows.IP_TOS, value)
	ipv6Err = windows.SetsockoptInt(windows.Handle(fd), ipProtoIPV6, ipv6TClass, value)
	return ipv4Err, ipv6Err
}