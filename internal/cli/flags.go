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

	// --target, -s
	target := flag.String("target", "", "--target [target]")
	flag.StringVar(target, "tar", "", "-tar [target]")

	// --ports, -p
	ports := flag.String("ports", "", "--ports [ports]")
	flag.StringVar(ports, "p", "", "-p [ports]")

	// --timeout, -t
	timeout := flag.Int("timeout", 1000, "--timeout [num]")
	flag.IntVar(timeout, "t", 1000, "-t [num]")

	// --json, -j
	asJSON := flag.String("json", "", "--json [path]")
	flag.StringVar(asJSON, "j", "", "-j")

	// --yaml, -y
	asYAML := flag.String("yaml", "", "--yaml [path]")
	flag.StringVar(asYAML, "y", "", "-y [path]")

	// --csv, -c
	asCSV := flag.String("csv", "", "--csv [path]")
	flag.StringVar(asCSV, "c", "", "-c [path]")

	// --parallel-hosts, -ph
	parallelHosts := flag.Int("parallel-hosts", 254, "--parallel-hosts [num]")
	flag.IntVar(parallelHosts, "ph", 254, "-ph [num]")

	// --parallel-ports, -pp
	parallelPorts := flag.Int("parallel-ports", 8, "--parallel-ports [num]")
	flag.IntVar(parallelPorts, "pp", 8, "-pp [num]")

	// Parse flags.
	flag.Parse()
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
