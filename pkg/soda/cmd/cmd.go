package cmd

import (
	"code.pikelabs.net/go/pkg/soda/cmd/mockbuild"
	"github.com/spf13/cobra"
)

func NewSodaCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "soda",
		Short: "RPM packaging tool",
	}

	cmd.AddCommand(mockbuild.NewMockbuildCmd())
	return cmd
}
