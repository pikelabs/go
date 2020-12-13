package build

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"code.pikelabs.net/go/pkg/soda"
	"code.pikelabs.net/go/pkg/soda/rpmbuild"
)

type Options struct {
	rpmbuildExtraArgs string
	buildRoot         string
}

func NewBuildCmd() *cobra.Command {
	var o Options

	cmd := &cobra.Command{
		Use:   "build",
		Short: "build SRPM from source code in using rpmbuild",
		Run: func(cmd *cobra.Command, args []string) {
			if err := o.Run(); err != nil {
				fmt.Fprintf(os.Stderr, "err: %s\n", err)
			}
		},
	}

	cmd.PersistentFlags().StringVar(&o.rpmbuildExtraArgs, "rpmbuild-extra", "", "apply extra arguments to rpmbuild")
	cmd.PersistentFlags().StringVar(&o.buildRoot, "buildroot", "", "set buildroot path")
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

	if err := rpmbuild.MkBuildroot(o.buildRoot, pkg.Specfile, pkg.Sources); err != nil {
		return err
	}

	return soda.RunCmd("rpmbuild",
		"--buildroot",
		o.buildRoot,
		o.rpmbuildExtraArgs,
		"-bs",
		pkg.Specfile)
}