package main

import (
	"os"
)

/**
 * soda is tool to work with source-git
 */
func main() {
	root := NewSodaCommand()
	if err := root.Execute(); err != nil {
		os.Exit(1)
	}
}
