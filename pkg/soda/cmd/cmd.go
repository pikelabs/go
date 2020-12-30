package cmd

import (
	"code.pikelabs.net/go/pkg/soda/cmd/build"
	"code.pikelabs.net/go/pkg/soda/cmd/mockbuild"
	"code.pikelabs.net/go/pkg/soda/cmd/prep"
	"code.pikelabs.net/go/pkg/soda/cmd/initialize"
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
