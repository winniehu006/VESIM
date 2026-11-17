package phone

import (
	"net"
	"strconv"
	"testing"
)

func TestSameUDPPort(t *testing.T) {
	cases := []struct {
		name       string
		configured string
		bound      string
		want       bool
	}{
		{"dual-stack form of the same port", ":7580", "[::]:7580", true},
		{"v4 form of the same port", "127.0.0.1:7580", "127.0.0.1:7580", true},
		{"fallback port", ":7580", "[::]:7581", false},
		{"ephemeral request", "127.0.0.1:0", "127.0.0.1:53412", true},
		{"unparsable configured address", "7580", "7580", true},
	}
	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			if got := sameUDPPort(testCase.configured, testCase.bound); got != testCase.want {
				t.Fatalf("sameUDPPort(%q, %q) = %v, want %v", testCase.configured, testCase.bound, got, testCase.want)
			}
		})
	}
}

func TestListenWebRTCUDPFallsBackWhenPortIsBusy(t *testing.T) {
	busy, err := net.ListenUDP("udp", &net.UDPAddr{IP: net.ParseIP("127.0.0.1")})
	if err != nil {
		t.Fatalf("listen on a busy probe socket: %v", err)
	}
	defer func() { _ = busy.Close() }()
	busyPort := busy.LocalAddr().(*net.UDPAddr).Port

	connection, boundAddress, err := listenWebRTCUDP("127.0.0.1:" + strconv.Itoa(busyPort))
	if err != nil {
		t.Fatalf("listenWebRTCUDP() error = %v", err)
	}
	defer func() { _ = connection.Close() }()

	_, boundPort, err := net.SplitHostPort(boundAddress)
	if err != nil {
		t.Fatalf("bound address %q is not host:port: %v", boundAddress, err)
	}
	if boundPort == strconv.Itoa(busyPort) {
		t.Fatalf("listenWebRTCUDP() kept the busy port %s", boundPort)
	}
	if sameUDPPort("127.0.0.1:"+strconv.Itoa(busyPort), boundAddress) {
		t.Fatalf("sameUDPPort() reported the fallback %s as the requested port", boundAddress)
	}
}

func TestListenWebRTCUDPUsesRequestedPort(t *testing.T) {
	// Ask the kernel for a free port, release it, and expect it back.
	probe, err := net.ListenUDP("udp", &net.UDPAddr{IP: net.ParseIP("127.0.0.1")})
	if err != nil {
		t.Fatalf("listen on a probe socket: %v", err)
	}
	port := probe.LocalAddr().(*net.UDPAddr).Port
	_ = probe.Close()

	address := "127.0.0.1:" + strconv.Itoa(port)
	connection, boundAddress, err := listenWebRTCUDP(address)
	if err != nil {
		t.Fatalf("listenWebRTCUDP(%q) error = %v", address, err)
	}
	defer func() { _ = connection.Close() }()

	if !sameUDPPort(address, boundAddress) {
		t.Fatalf("listenWebRTCUDP(%q) bound %q", address, boundAddress)
	}
}
