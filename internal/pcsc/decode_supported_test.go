//go:build darwin || linux

package pcsc

import "testing"

// decodeMultiString only exists on the platforms with a real PC/SC backend.
func TestDecodeMultiString(t *testing.T) {
	want := []string{"Reader A", "Reader B"}
	got := decodeMultiString([]byte("Reader A\x00Reader B\x00\x00"))
	if len(got) != len(want) || got[0] != want[0] || got[1] != want[1] {
		t.Fatalf("decodeMultiString = %#v", got)
	}
}
