package prep

import (
	"fmt"
	"errors"
	"io/ioutil"
	"github.com/spf13/cobra"
	"code.pikelabs.net/go/pkg/soda"
	"code.pikelabs.net/go/pkg/soda/rpmbuild"

)

type Options struct {
	File string
	BuildRoot string
}

func NewPrepCmd() *cobra.Command {
	var opt Options
	c := &cobra.Command {
		Use: "prep",
		Short: "Run prep stage of RPM build",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("command requires exectly one argument")
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			if err := opt.Prep(cmd, args); err != nil {
				return
			}
			opt.Run()
		},
	}

	c.Flags().StringVarP(&opt.BuildRoot, "buildroot", "C", "", "Base directory for rpm build")
	return c
}


func (opts *Options) Prep(cmd *cobra.Command, args []string) error {
	if len(args) < 1 {
		return errors.New("command requires exectly one argument")
	}
	opts.File = args[0]
	if len(opts.BuildRoot) < 1 {
		dir, err := ioutil.TempDir("", "soda")
		if err != nil {
			return err
		}
		opts.BuildRoot = dir
	}
	return nil
}


func (opts *Options) Run() {
	rpmbuild.MkBuildroot(opts.BuildRoot, opts.File, []string{"./"})
	soda.RunCmd("rpmbuild", "--define", fmt.Sprintf("_topdir %s", opts.BuildRoot), "-bp", opts.File)
}
