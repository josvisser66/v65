package main

import (
	"flag"
	"fmt"
	"os"
	"v65/asm"
)

func main() {
	flag.Parse()

	for _, sourceFile := range flag.Args() {
		if err := asm.Assemble(sourceFile); err != nil {
			fmt.Fprintf(os.Stderr, "assembly error: %v\n", err)
		}
	}
}
