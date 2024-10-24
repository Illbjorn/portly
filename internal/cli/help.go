package cli

import (
	"fmt"
	"os"

	"github.com/illbjorn/fstr"
)

var println = fmt.Println

var helpText = fstr.Pairs(`
{magenta}Portly, a simple, stalwart port scanner.{reset}

{cyan}Docs{reset} : https://github.com/Illbjorn/portly
{cyan}Bugs{reset} : https://github.com/Illbjorn/portly/issues

To perform a basic network scan:

  {green}portly -s 192.168.255.0/24 -p 80,443{reset}

Options:
  --target,         -tar  The subnet to be scanned. This can be provided in CIDR
                          format (192.168.255.0/24) or as a single IP
                          (192.168.255.1).
  --ports,          -p    The ports to scan. These may be multiple ports,
                          comma-delimited.
  --timeout,        -t    The time to wait (in milliseconds) for a response from
                          the current target's port.
  --parallel-hosts, -ph   The number of hosts to scan concurrently.
  --parallel-ports, -pp   The number of ports per-host to scan concurrently.
  --help,           -h    How you got here!
`,
	"cyan", "\033[0m",
	"yellow", "\033[33m",
	"magenta", "\033[35m",
	"reset", "\033[0m",
	"green", "\033[32m",
)

func helpAndExit(code int) {
	println(helpText)
	os.Exit(code)
}
