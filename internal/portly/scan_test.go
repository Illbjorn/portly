package portly_test

import (
	"net/netip"
	"testing"

	"github.com/illbjorn/portly/internal/portly"
)

var network, _ = netip.ParsePrefix("192.168.255.0/24")
var ports = []int{21, 22, 25, 80, 443}
var target = portly.NewTarget(network)

func BenchmarkScan(b *testing.B) {
	for range b.N {
		portly.Scan(target, ports...)
	}
}
