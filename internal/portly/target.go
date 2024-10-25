package portly

import "net/netip"

func NewTarget(targets ...netip.Prefix) Target {
	return Target{Targets: targets}
}

type Target struct {
	Targets []netip.Prefix
}

func (t *Target) Gen() chan netip.Addr {
	ch := make(chan netip.Addr)

	go func() {
		defer close(ch)

		for _, target := range t.Targets {
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
