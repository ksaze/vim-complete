package main

import (
	"fmt"
	"os"

	"github.com/ksaze/vim-complete/commands"
)

func main() {
	if err := commands.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, "vc:", err)
		os.Exit(1)
	}
}
