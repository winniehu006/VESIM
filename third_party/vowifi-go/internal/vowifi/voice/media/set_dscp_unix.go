//go:build !windows

package media

import "syscall"

// setTrafficClass applies the IPv4 TOS value and the IPv6 traffic class value.
func setTrafficClass(fd uintptr, value int) (ipv4Err, ipv6Err error) {
	ipv4Err = syscall.SetsockoptInt(int(fd), syscall.IPPROTO_IP, syscall.IP_TOS, value)
	ipv6Err = syscall.SetsockoptInt(int(fd), syscall.IPPROTO_IPV6, syscall.IPV6_TCLASS, value)
	return ipv4Err, ipv6Err
}