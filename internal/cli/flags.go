package cli

import (
	"flag"
	"fmt"
	"os"
	"time"

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
	Timeout       time.Duration
}

func parseFlags() flags {
	fs := flag.NewFlagSet("", flag.ExitOnError)
	f := flags{}

	// --target, -t
	// Target to scan.
	target := fs.String("target", "", "--target [target]")
	fs.StringVar(target, "t", "", "-t [target]")

	// --ports, -p
	// Ports to scan.
	ports := fs.String("ports", "", "--ports [ports]")
	fs.StringVar(ports, "p", "", "-p [ports]")

	// --timeout, -to
	// Time to wait, per-port per-host, before considering it "closed".
	timeout := fs.Duration("timeout", 1000*time.Millisecond, "--timeout [duration]")
	fs.DurationVar(timeout, "to", 1000*time.Millisecond, "-to [duration]")

	// --json, -j
	// Serialize the results as JSON and write to disk.
	asJSON := fs.String("json", "", "--json [path]")
	fs.StringVar(asJSON, "j", "", "-j")

	// --yaml, -y
	// Serialize the results as YAML and write to disk.
	asYAML := fs.String("yaml", "", "--yaml [path]")
	fs.StringVar(asYAML, "y", "", "-y [path]")

	// --csv, -c
	// Serialize the results as CSV and write to disk.
	asCSV := fs.String("csv", "", "--csv [path]")
	fs.StringVar(asCSV, "c", "", "-c [path]")

	// --parallel-hosts, -ph
	// Number of hosts to scan in parallel.
	parallelHosts := fs.Int("parallel-hosts", 254, "--parallel-hosts [num]")
	fs.IntVar(parallelHosts, "ph", 254, "-ph [num]")

	// --parallel-ports, -pp
	// Number of ports per-host to scan in parallel.
	parallelPorts := fs.Int("parallel-ports", 8, "--parallel-ports [num]")
	fs.IntVar(parallelPorts, "pp", 8, "-pp [num]")

	// Parse flags.
	fs.Parse(os.Args[1:])

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
