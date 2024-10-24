package scan

import (
	"net"
	"net/netip"
	"sort"
	"strconv"
	"time"

	"golang.org/x/sync/errgroup"
)

const (
	PORT_STATUS_CLOSED = "CLOSED"
	PORT_STATUS_OPEN   = "OPEN"
)

var (
	Timeout = 1000 * time.Millisecond
	itoa    = strconv.Itoa

	ConcurrentHostScans = 254
	ConcurrentPortScans = 8
)

func parallelHostScan(network netip.Prefix, ports ...int) chan HostResult {
	// Init the firehose!
	ch := make(chan HostResult, ConcurrentHostScans)

	go func() {
		defer close(ch)

		// Init the error group, limiting concurrency to the parallel host scan
		// limit.
		g := errgroup.Group{}
		g.SetLimit(ConcurrentHostScans)

		// Prepare the first address to scan.
		// We advance immediately as to skip the network ID.
		next := network.Addr()
		for {
			// We look position+1 here to skip the network broadcast.
			if !network.Contains(next) {
				break
			}

			// Create another reference to `next` to close over.
			ip := next

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

			// Advance to the next address to scan.
			next = next.Next()
		}

		// Wait for all scanners to complete.
		g.Wait()
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

		g.Wait()
	}()

	return ch
}

func Host(prefix netip.Prefix, ports ...int) ScanResult {
	res := newScanResult(prefix, ports...)
	host := newHostResult(prefix.Addr())

	// Iterate all ports to be scanned.
	for port := range parallelPortScan(prefix.Addr(), ports...) {
		host.Ports = append(host.Ports, port)
	}

	// Sort the ports slice.
	sort.Slice(host.Ports, func(i, j int) bool {
		return host.Ports[i].Port < host.Ports[j].Port
	})

	res.Hosts = []HostResult{host}

	return res
}

func Range(network netip.Prefix, ports ...int) ScanResult {
	res := newScanResult(network, ports...)

	for scan := range parallelHostScan(network, ports...) {
		// Sort each host result port slice by port number.
		sort.Slice(scan.Ports, func(i, j int) bool {
			return scan.Ports[i].Port < scan.Ports[j].Port
		})

		// Add the result.
		res.Hosts = append(res.Hosts, scan)
	}

	// Sort the host results slice by host IP.
	sort.Slice(res.Hosts, func(i, j int) bool {
		return res.Hosts[i].Host.Compare(res.Hosts[j].Host) < 0
	})

	return res
}
