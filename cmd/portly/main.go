package main

import (
	"os"
	"time"

	"github.com/illbjorn/portly/internal/cli"
)

func main() {
	// Measure runtime.
	start := time.Now()
	cli.Run(os.Args)
	runtime := time.Since(start)
	println("---")
	println("Runtime : " + runtime.String())
	println("---")
}
