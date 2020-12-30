package main

import (
	"code.pikelabs.net/go/soda/cmd/build"
	"code.pikelabs.net/go/soda/cmd/initialize"
	"code.pikelabs.net/go/soda/cmd/mockbuild"
	"code.pikelabs.net/go/soda/cmd/prep"
	"github.com/spf13/cobra"
)

func NewSodaCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "soda",
		Short: "RPM packaging tool",
	}

	cmd.AddCommand(build.NewBuildCmd())
	cmd.AddCommand(mockbuild.NewMockbuildCmd())
	cmd.AddCommand(prep.NewPrepCmd())
	cmd.AddCommand(initialize.NewInitCmd())
	return cmd
}
