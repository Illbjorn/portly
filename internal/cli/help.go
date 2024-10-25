package cli

import (
	"fmt"
	"os"

	"github.com/illbjorn/fstr"
)

var println = fmt.Println

var helpText = fstr.Pairs(`
{magenta}Portly, your stalwart port scanning companion.{reset}

{cyan}Docs{reset} : https://github.com/Illbjorn/portly
{cyan}Bugs{reset} : https://github.com/Illbjorn/portly/issues

To perform a basic network scan:

  {green}portly -t 192.168.255.0/24 -p 80,443{reset}

Options:
  --target,         -t    The subnet to be scanned. This can be provided in CIDR
                          format (192.168.255.0/24), as a single IP
                          (192.168.255.1) or a hostname.
  --ports,          -p    The ports to scan. This value may be multiple ports,
                          comma-delimited.
  --timeout,        -to   The time to wait for a response from the current
                          target's port. This must be a valid Go Duration.
                          Examples: 1000ms, 1s, 1m30s.
  --parallel-hosts, -ph   The number of hosts to scan concurrently.
  --parallel-ports, -pp   The number of ports per-host to scan concurrently.
  --json,           -j    Serialize the result as JSON and write to disk.
  --yaml,           -y    Serialize the result as YAML and write to disk.
  --csv,            -c    Serialize the result as CSV and write to disk.
  --help,           -h    How you got here!

`,
	"cyan", "\033[0m",
	"yellow", "\033[33m",
	"magenta", "\033[35m",
	"reset", "\033[0m",
	"green", "\033[32m",
)

func helpAndExit(code int) {
	_, _ = println(helpText)
	os.Exit(code)
}
