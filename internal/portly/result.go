package portly

import "net/netip"

func newScanResult(target Target, ports ...int) Result {
	return Result{
		Target: target,
		Ports:  ports,
	}
}

type Result struct {
	Target Target
	Ports  []int
	Hosts  []HostResult
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
