package cli

import (
	"flag"

	"github.com/illbjorn/portly/internal/assert"
)

type flags struct {
	Target        string
	Ports         string
	AsJSON        string
	AsYAML        string
	AsCSV         string
	ParallelHosts int
	ParallelPorts int
	Timeout       int
}

func parseFlags() flags {
	f := flags{}

	// --target, -t
	// Target to scan.
	target := flag.String("target", "", "--target [target]")
	flag.StringVar(target, "t", "", "-t [target]")

	// --ports, -p
	// Ports to scan.
	ports := flag.String("ports", "", "--ports [ports]")
	flag.StringVar(ports, "p", "", "-p [ports]")

	// --timeout, -to
	// Time to wait, per-port per-host, before considering it "closed".
	timeout := flag.Int("timeout", 1000, "--timeout [num]")
	flag.IntVar(timeout, "to", 1000, "-to [num]")

	// --json, -j
	// Serialize the results as JSON and write to disk.
	asJSON := flag.String("json", "", "--json [path]")
	flag.StringVar(asJSON, "j", "", "-j")

	// --yaml, -y
	// Serialize the results as YAML and write to disk.
	asYAML := flag.String("yaml", "", "--yaml [path]")
	flag.StringVar(asYAML, "y", "", "-y [path]")

	// --csv, -c
	// Serialize the results as CSV and write to disk.
	asCSV := flag.String("csv", "", "--csv [path]")
	flag.StringVar(asCSV, "c", "", "-c [path]")

	// --parallel-hosts, -ph
	// Number of hosts to scan in parallel.
	parallelHosts := flag.Int("parallel-hosts", 254, "--parallel-hosts [num]")
	flag.IntVar(parallelHosts, "ph", 254, "-ph [num]")

	// --parallel-ports, -pp
	// Number of ports per-host to scan in parallel.
	parallelPorts := flag.Int("parallel-ports", 8, "--parallel-ports [num]")
	flag.IntVar(parallelPorts, "pp", 8, "-pp [num]")

	// Parse flags.
	flag.Parse()

	// Assign flag values.
	f.Target = *target
	f.Ports = *ports
	f.Timeout = *timeout
	f.AsJSON = *asJSON
	f.AsYAML = *asYAML
	f.AsCSV = *asCSV
	f.ParallelHosts = *parallelHosts
	f.ParallelPorts = *parallelPorts

	// Confirm we have required values.
	assert.NE(f.Target, "", "The --target, -s flag is required.")
	assert.NE(f.Ports, "", "The --prefix, -p flag is required.")

	return f
}
