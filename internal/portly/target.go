package portly

import (
	"encoding/json"
	"net/netip"
)

func NewTarget(targets ...netip.Prefix) Target {
	return Target{targets: targets}
}

type Target struct {
	targets []netip.Prefix
}

var _ json.Marshaler = Target{}

// We implement the json.Marshaler interface to avoid exporting field `targets`.
//
// Unclear if this actually makes sense \o/.
//
// NOTE: Currently, we write the byte-slice-converted address string and prefix
// length with a '/' in the middle separately. While `netip.Prefix` does have
// `fmt.Stringer` support this separation is intentional as I will eventually
// come back to this and write IPv4/6 address bytes directly to the `out` slice
// to skip the string allocation.
func (t Target) MarshalJSON() ([]byte, error) {
	var out []byte

	// Open the array.
	out = append(out, '[')
	for i, target := range t.targets {
		// Open the array string member.
		out = append(out, '"')

		// Marshal the address.
		addr := target.Addr().String()
		out = append(out, []byte(addr)...) // TODO: This is the lazy way out. Need to clean this up.

		// Marshal the CIDR block prefix length.
		out = append(out, '/')
		prefixLen := itoa(target.Bits())
		out = append(out, []byte(prefixLen)...) // TODO: This is the lazy way out. Need to clean this up.

		// End the string array member.
		out = append(out, '"')

		// Append a comma to all but the last in the series.
		if i < len(t.targets)-1 {
			out = append(out, ',')
		}
	}
	// Close the array.
	out = append(out, ']')

	return out, nil
}

// A simple generator, sending all targets across a channel. This sends both
// individual addresses and CIDR blocks.
func (t *Target) gen() chan netip.Addr {
	ch := make(chan netip.Addr)

	go func() {
		defer close(ch)

		for _, target := range t.targets {
			next := target.Addr()
			for {
				if !target.Contains(next) {
					break
				}

				ch <- next

				next = next.Next()
			}
		}
	}()

	return ch
}
