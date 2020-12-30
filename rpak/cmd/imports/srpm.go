package imports

import (
	"fmt"
	"os"

	"code.pikelabs.net/go/rpak"
	"github.com/spf13/cobra"
)

func fatal(format string, a ...interface{}) {
	fmt.Fprintf(os.Stderr, format, a...)
	os.Exit(1)
}

func NewImportSRPMCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "import-srpm",
		Short: "Import SRPM to repo",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) < 1 {
				fatal("err: command requires one argument\n")
			}
			cwd, err := os.Getwd()
			if err != nil {
				fatal("err: unable to get CWD: %s", err)
			}
			w := rpak.NewWorkflow(cwd)
			w.LoadRepo()
			err = w.ImportSRPM(os.Args[1])
			if err != nil {
				fatal("err: import failed %s\n", err.Error())
			}
		},
	}

	return cmd
}
