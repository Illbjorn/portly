# Overview

Welcome to Portly, your stalwart port scanning companion!

Portly is a simple, blazingly fast port scanner capable of: single IP, CIDR
block and hostname (`A` records only, for now) targeting.

Portly implements two major performance boons:
- TCP connectivity is performed via the Go stdlib's `net.DialTimeout` which
allows us to skip the `ACK`+ portion of the TCP handshake.
- All host and port scanning is concurrent-by-default!
  - By default, Portly will scan up to 254 addresses at a time, and will check 8-ports *per-host* at a time.
  - That's a *default* of 2032 concurrent operations for any given scan.
  - However, considering the efficiency of both the Go scheduler and the `net.Dialer`, your system resources will barely budge.

# Quickstart

To begin using Portly, simply `go install` it:

```sh
go install github.com/illbjorn/portly/cmd/portly@main
```

To perform a basic network scan:
```sh
portly -t 192.168.255.0/24 -p 80,443
```

# Usage

```sh
Portly, your stalwart port scanning companion.

Docs : https://github.com/Illbjorn/portly
Bugs : https://github.com/Illbjorn/portly/issues

To perform a basic network scan:

  portly -t 192.168.255.0/24 -p 80,443

Options:
  --target,         -t    The subnet to be scanned. This can be provided in CIDR
                          format (192.168.255.0/24), as a single IP
                          (192.168.255.1) or a hostname.
  --ports,          -p    The ports to scan. This value may be multiple ports,
                          comma-delimited.
  --timeout,        -to   The time to wait for a response from the current
                          target's port. This must be a valid Go Duration.
                          Examples: 1000ms, 1s, 1m30s.
  --open-only,      -oo   Filters results to only those where the port was
                          'open'.
  --parallel-hosts, -ph   The number of hosts to scan concurrently.
  --parallel-ports, -pp   The number of ports per-host to scan concurrently.
  --json,           -j    Serialize the result as JSON and write to disk.
  --yaml,           -y    Serialize the result as YAML and write to disk.
  --csv,            -c    Serialize the res
```
