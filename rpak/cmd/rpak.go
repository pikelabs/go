package main

import (
	"code.pikelabs.net/go/rpak/cmd/imports"
)

func main() {
	cmd := imports.NewImportSRPMCmd()
	if err := cmd.Execute(); err != nil {
		panic(err)
	}
}
