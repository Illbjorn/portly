package cli

import (
	"encoding/csv"
	"encoding/json"
	"os"

	"github.com/illbjorn/portly/internal/assert"
	"github.com/illbjorn/portly/internal/portly"
	"gopkg.in/yaml.v3"
)

func writeResultJSON(results portly.Result, path string) {
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

func writeResultYAML(results portly.Result, path string) {
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

func writeResultCSV(results portly.Result, path string) {
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
