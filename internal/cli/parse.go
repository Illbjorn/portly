package cli

import (
	"net"
	"net/netip"
	"regexp"
	"strings"

	"github.com/illbjorn/portly/internal/assert"
	"github.com/illbjorn/portly/internal/portly"
)

const (
	// IPv4
	psIsIPv4IP      = `^ *[0-9]+\.[0-9]+\.[0-9]+\.[0-9]+ *$`
	psIsIPv4Network = `^ *[0-9]+\.[0-9]+\.[0-9]+\.[0-9]+/(?:0|[1-9][0-9]*) *$`
	// IPv6
	psIsIPv6IP      = `^ *(?:[A-Fa-f0-9]{1,4}:){7}[A-Fa-f0-9]{1,4} *$`
	psIsIPv6Network = `^ *[A-Fa-f0-9]{1,4}:[A-Fa-f0-9]{1,4}::/[0-9]{1,3}`
	// Hostname
	psIsHostname = `^ *[a-zA-Z0-9][a-zA-Z.0-9\-]+[a-zA-Z0-9] *$`
)

var (
	// IPv4
	pIsIPv4IP      = regexp.MustCompile(psIsIPv4IP)
	pIsIPv4Network = regexp.MustCompile(psIsIPv4Network)
	// IPv6
	pIsIPv6IP      = regexp.MustCompile(psIsIPv6IP)
	pIsIPv6Network = regexp.MustCompile(psIsIPv6Network)
	// Hostname
	pIsHostname = regexp.MustCompile(psIsHostname)
)

func parseTarget(target string) portly.Target {
	// Identify the target type.
	switch {
	// Target is a single IPv4 IP.
	case pIsIPv4IP.MatchString(target):
		return parseIPv4IP(target)

		// Target is a single IPv6 IP.
	case pIsIPv6IP.MatchString(target):
		return parseIPv6IP(target)

	// Target is an IPv4/IPv6 CIDR network.
	case pIsIPv4Network.MatchString(target), pIsIPv6Network.MatchString(target):
		return parseNetwork(target)

	// Target is a hostname.
	case pIsHostname.MatchString(target):
		return parseHostname(target)

	// Unrecognized target.
	default:
		println("Failed to classify target: '" + target + "'.")
		println("Only a single IP, a single CIDR network and hostnames are supported currently.")
		exit(1)

		return portly.Target{}
	}
}

func parseIPv4IP(target string) portly.Target {
	// Parse the single IP as a /32 network.
	prefix, err := netip.ParsePrefix(target + "/32")
	assert.EQ(err, nil, "Failed to parse '"+target+"' as a /32 IPv4 IP.")

	return portly.NewTarget(prefix)
}

func parseIPv6IP(target string) portly.Target {
	// Parse the single IP as a /128 network.
	prefix, err := netip.ParsePrefix(target + "/128")
	assert.EQ(err, nil, "Failed to parse '"+target+"' as a /128 IPv6 IP.")

	return portly.NewTarget(prefix)
}

func parseNetwork(target string) portly.Target {
	// Parse the network as a CIDR network.
	prefix, err := netip.ParsePrefix(target)
	assert.EQ(err, nil, "Failed to parse '"+target+"' as a CIDR network.")

	return portly.NewTarget(prefix)
}

func parseHostname(target string) portly.Target {
	var networks []netip.Prefix

	// Resolve the hostname.
	ips, err := net.LookupHost(target)
	assert.EQ(err, nil, "Failed to lookup hostname: '"+target+"'.")

	// Create a /32 or /128 netip.Prefix from each IP.
	var prefix netip.Prefix
	for _, ip := range ips {
		if pIsIPv4IP.MatchString(ip) {
			prefix, err = netip.ParsePrefix(ip + "/32")
		} else {
			prefix, err = netip.ParsePrefix(ip + "/128")
		}
		assert.EQ(err, nil, "Failed to parse resolved IP '"+ip+"' as a single IP.")
		networks = append(networks, prefix)
	}

	return portly.NewTarget(networks...)
}

func parsePorts(portStr string) ([]int, int) {
	// Split the string on comma.
	portStrs := strings.Split(portStr, ",")

	// Process each port string.
	var longest int
	var ports []int
	for _, port := range portStrs {
		port = strings.TrimSpace(port)

		// This is just for result output purposes at the end.
		if len(port) > longest {
			longest = len(port)
		}

		ports = append(ports, cint(port))
	}

	// We must have at least one port.
	assert.GT(len(ports), 0, "At least one port to scan is required.")

	return ports, longest
}
