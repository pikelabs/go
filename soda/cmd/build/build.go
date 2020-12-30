package build

import (
	"encoding/json"
	"fmt"
	"os"

	"io/ioutil"

	"github.com/spf13/cobra"

	"code.pikelabs.net/go/soda"
	"code.pikelabs.net/go/soda/rpmbuild"
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
			if err := o.Prep(cmd, args); err != nil {
				fmt.Fprintf(os.Stderr, "err: %s\n", err)
				return
			}
			if err := o.Run(); err != nil {
				fmt.Fprintf(os.Stderr, "err: %s\n", err)
			}
		},
	}

	cmd.PersistentFlags().StringVar(&o.rpmbuildExtraArgs, "rpmbuild-extra", "", "apply extra arguments to rpmbuild")
	cmd.PersistentFlags().StringVar(&o.buildRoot, "buildroot", "", "set buildroot path")
	return cmd
}

func (o *Options) Prep(cmd *cobra.Command, args []string) error {
	if len(o.buildRoot) < 1 {
		dir, err := ioutil.TempDir("", "soda")
		if err != nil {
			return err
		}
		o.buildRoot = dir
	}
	return nil
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

	cwd, err := os.Getwd()
	if err != nil {
		return err
	}
	return soda.RunCmd("rpmbuild",
		"-bs",
		"--define",
		fmt.Sprintf("_topdir %s", o.buildRoot),
		"--define",
		fmt.Sprintf("_srcrpmdir %s", cwd),
		pkg.Specfile)
}
