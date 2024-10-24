package scan_test

import (
	"net/netip"
	"testing"

	"github.com/illbjorn/portly/internal/scan"
)

var network, _ = netip.ParsePrefix("192.168.255.0/24")
var ports = []int{21, 22, 25, 80, 443}

func BenchmarkRange(b *testing.B) {
	for range b.N {
		scan.Range(network, ports...)
	}
}
