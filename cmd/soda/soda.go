package main

import (
	"os"

	"code.pikelabs.net/go/pkg/soda/cmd"
)

/**
 * soda is tool to work with source-git
 */
func main() {
	root := cmd.NewSodaCommand()
	if err := root.Execute(); err != nil {
		os.Exit(1)
	}
}
