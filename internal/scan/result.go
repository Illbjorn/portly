package scan

import "net/netip"

func newScanResult(addr netip.Prefix, ports ...int) ScanResult {
	return ScanResult{
		Network: addr,
		Ports:   ports,
	}
}

type ScanResult struct {
	Network netip.Prefix
	Ports   []int
	Hosts   []HostResult
}

func newHostResult(addr netip.Addr) HostResult {
	return HostResult{Host: addr}
}

type HostResult struct {
	Host  netip.Addr
	Ports []PortResult
}

func newPortResult(port int, status string) PortResult {
	return PortResult{Port: port, Status: status}
}

type PortResult struct {
	Status string
	Port   int
}
