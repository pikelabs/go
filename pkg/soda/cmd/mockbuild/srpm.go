package mockbuild

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"code.pikelabs.net/go/pkg/soda"
)

type Options struct {
	extraMockArg string
	buildRoot    string
}

func NewMockbuildCmd() *cobra.Command {
	var o Options

	cmd := &cobra.Command{
		Use:   "mockbuild",
		Short: "build SRPM from source code in using mock",
		Run: func(cmd *cobra.Command, args []string) {
			if err := o.Run(); err != nil {
				fmt.Fprintf(os.Stderr, "err: %s\n", err)
			}
		},
	}

	cmd.PersistentFlags().StringVar(&o.extraMockArg, "mock-extra", "", "apply extra arguments to mock")
	return cmd
}

func (o Options) Run() error {
	f, err := os.Open(".sodafile")
	if err != nil {
		return err
	}
	defer f.Close()
	var pkg soda.Package
	err = json.NewDecoder(f).Decode(&pkg)
	if err != nil {
		return err
	}

	for i, step := range pkg.PrebuildSteps {
		if err1 := step.Run(); err1 != nil {
			err = fmt.Errorf("cmd(%d): failed: %s", i, err1.Error())
			return err
		}
	}

	mockCmd := fmt.Sprintf("mock --buildsrpm --spec=%s --sources=%s %s", pkg.Specfile, strings.Join(pkg.Sources, ","), o.extraMockArg)
	return soda.Shellout(mockCmd)
}
