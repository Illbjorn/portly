package portly

import (
	"net"
	"net/netip"
	"slices"
	"time"

	"golang.org/x/sync/errgroup"
)

const (
	PORT_STATUS_CLOSED = "CLOSED"
	PORT_STATUS_OPEN   = "OPEN"
)

var (
	// Default operating values.
	Timeout             = 1000 * time.Millisecond
	ConcurrentHostScans = 254
	ConcurrentPortScans = 8
)

func parallelHostScan(target Target, ports ...int) chan HostResult {
	// Init the firehose!
	ch := make(chan HostResult, ConcurrentHostScans)

	go func() {
		defer close(ch)

		// Init the error group, limiting concurrency to the parallel host scan
		// limit.
		g := errgroup.Group{}
		g.SetLimit(ConcurrentHostScans)

		// Scan all target IPs.
		for ip := range target.gen() {
			// Spawn the scanner.
			// This respects `ConcurrentHostScans` as a concurrency limit.
			g.Go(
				func() error {
					// Init the host scan result.
					scan := newHostResult(ip)

					// Perform the port scan.
					for portres := range parallelPortScan(ip, ports...) {
						scan.Ports = append(scan.Ports, portres)
					}

					// Send the result.
					ch <- scan

					return nil
				})
		}

		// Wait for all scanners to complete.
		_ = g.Wait() // TODO: Handle.
	}()

	return ch
}

func parallelPortScan(addr netip.Addr, ports ...int) chan PortResult {
	ch := make(chan PortResult, ConcurrentPortScans)

	go func() {
		defer close(ch)

		g := errgroup.Group{}
		g.SetLimit(ConcurrentPortScans)

		for _, port := range ports {
			g.Go(
				func() error {
					// Attempt a TCP connection.
					conn, err := net.DialTimeout("tcp", addr.String()+":"+itoa(port), Timeout)
					if err != nil {
						// Indicate failure.
						ch <- newPortResult(port, PORT_STATUS_CLOSED)
						return nil
					}
					defer conn.Close()

					// Indicate success.
					ch <- newPortResult(port, PORT_STATUS_OPEN)
					return nil
				})
		}

		_ = g.Wait() // TODO: Handle.
	}()

	return ch
}

func Scan(target Target, ports ...int) Result {
	res := newScanResult(target, ports...)

	for scan := range parallelHostScan(target, ports...) {
		// Sort each host result port slice by port number.
		slices.SortFunc(scan.Ports, func(left, right PortResult) int {
			switch {
			// left > right
			case left.Port > right.Port:
				return 1

			// left < right
			case left.Port > right.Port:
				return -1

			// left == right
			default:
				return 0
			}
		})

		// Add the result.
		res.Hosts = append(res.Hosts, scan)
	}

	// Sort the host results slice by host IP.
	slices.SortFunc(res.Hosts, func(left, right HostResult) int {
		return left.Host.Compare(right.Host)
	})

	return res
}
