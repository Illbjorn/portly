package cli

import (
	"github.com/illbjorn/portly/internal/assert"
	"github.com/illbjorn/portly/internal/portly"
)

func Run(args []string) {
	// We expect > 1 arg.
	if len(args) < 2 {
		helpAndExit(1)
	}

	// Parse flags.
	flags := parseFlags()

	// Parse the target.
	target := parseTarget(flags.Target)

	// Parse the port(s).
	ports, longestPort := parsePorts(flags.Ports)

	// Assign desired timeout duration.
	portly.Timeout = flags.Timeout

	// Assign desired concurrency values.
	portly.ConcurrentHostScans = flags.ParallelHosts
	portly.ConcurrentPortScans = flags.ParallelPorts

	// Perform the scan.
	res := portly.Scan(target, ports...)
	assert.GT(len(res.Hosts), 0, "Received no scan results.")

	// Output to JSON if specified.
	if flags.AsJSON != "" {
		writeResultJSON(res, flags.AsJSON)
	}

	// Output to YAML if specified.
	if flags.AsYAML != "" {
		writeResultYAML(res, flags.AsYAML)
	}

	// Output to CSV if specified.
	if flags.AsCSV != "" {
		writeResultCSV(res, flags.AsCSV)
	}

	// Output results to stdout.
	for _, r := range res.Hosts {
		_, _ = println(r.Host)
		for _, port := range r.Ports {
			pstr := itoa(port.Port)
			pstatus := "  " + pstr + repeat(" ", longestPort+1-len(pstr)) + ": " + port.Status
			_, _ = println(pstatus)
		}
	}
}
