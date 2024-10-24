package cli

import (
	"encoding/csv"
	"encoding/json"
	"net/netip"
	"os"
	"strings"
	"time"

	"github.com/illbjorn/fstr"
	"github.com/illbjorn/portly/internal/assert"
	"github.com/illbjorn/portly/internal/scan"
	"gopkg.in/yaml.v3"
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
	ports := parsePorts(flags.Ports)

	// Assign desired timeout duration.
	scan.Timeout = time.Duration(flags.Timeout) * time.Millisecond

	// Assign desired concurrency values.
	scan.ConcurrentHostScans = flags.ParallelHosts
	scan.ConcurrentPortScans = flags.ParallelPorts

	// Perform the scan.
	res := scan.Range(target, ports...)
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
			_, _ = println("  " + itoa(port.Port) + ": " + port.Status)
		}
	}
}

func parseTarget(target string) netip.Prefix {
	var prefix netip.Prefix
	var err error

	// If the string doesn't contain a forward slash, parse as a single IP.
	if !strings.Contains(target, "/") {
		prefix, err = netip.ParsePrefix(target + "/32")
	} else {
		prefix, err = netip.ParsePrefix(target)
	}
	assert.EQ(err, nil, fstr.Pairs(
		"Failed to parse provided subnet: {subnet}.",
		"subnet", target,
	))

	return prefix
}

func parsePorts(portStr string) []int {
	// Split the string on comma.
	portStrs := strings.Split(portStr, ",")

	// Process each port string.
	var ports []int
	for _, port := range portStrs {
		port = strings.TrimSpace(port)
		ports = append(ports, cint(port))
	}

	// We must have at least one port.
	assert.GT(len(ports), 0, "At least one port to scan is required.")

	return ports
}

func writeResultJSON(results scan.ScanResult, path string) {
	// Open a writable file stream.
	f, err := os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0600)
	assert.EQ(err, nil, "Failed to open writable stream to JSON output file: "+path+".")
	defer f.Close()

	// Encode the results as JSON.
	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	err = enc.Encode(results)
	assert.EQ(err, nil, "Failed to encode results as JSON.")
}

func writeResultYAML(results scan.ScanResult, path string) {
	// Open a writable file stream.
	f, err := os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0600)
	assert.EQ(err, nil, "Failed to open writable stream to YAML output file: "+path+".")
	defer f.Close()

	// Encode the results as YAML.
	enc := yaml.NewEncoder(f)
	enc.SetIndent(2)
	err = enc.Encode(results)
	assert.EQ(err, nil, "Failed to encode results as YAML.")
}

func writeResultCSV(results scan.ScanResult, path string) {
	// Open a writable file stream.
	f, err := os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0600)
	assert.EQ(err, nil, "Failed to open writable stream to CSV output file: "+path+".")
	defer f.Close()

	// Encode the results as CSV.
	enc := csv.NewWriter(f)
	enc.UseCRLF = false

	// Write the header row.
	header := []string{"host"}
	for _, port := range results.Ports {
		header = append(header, itoa(port))
	}
	_ = enc.Write(header)
	enc.Flush()

	for _, result := range results.Hosts {
		row := make([]string, len(result.Ports)+1)
		row[0] = result.Host.String()
		for i, port := range result.Ports {
			row[i+1] = port.Status
		}
		_ = enc.Write(row)
		enc.Flush()
	}
}
