package cli

import (
	"os"
	"strconv"
	"strings"

	"github.com/illbjorn/fstr"
	"github.com/illbjorn/portly/internal/assert"
)

var itoa = strconv.Itoa
var repeat = strings.Repeat
var exit = os.Exit

func cint(d string) int {
	i, err := strconv.ParseInt(d, 10, 32)
	assert.EQ(err, nil, fstr.Pairs(
		"Failed to parse string as int64: {d}.",
		"d", d,
	))

	return int(i)
}
