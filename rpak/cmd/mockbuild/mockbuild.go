package mockbuild

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

type Options struct {
	buildroot string
	arch      string
	srpm      string
}

func NewMockbuildCmd() *cobra.Command {
	var opts Options
	command := &cobra.Command{
		Use:   "mockbuild [flags] srpm",
		Short: "Local build of SRPM using using mock",
		Run: func(cmd *cobra.Command, args []string) {
			if err := opts.Prepare(cmd, args); err != nil {
				fmt.Fprintf(os.Stderr, "err: %s\n", err.Error())
				return
			}
			opts.Run()
		},
	}
	command.Flags().StringVarP(&opts.buildroot, "root", "c", "", "Set mock  chroot config")
	command.Flags().StringVarP(&opts.buildroot, "arch", "a", "", "Set mock build architecture")
	return command
}

func (opts *Options) Prepare(cmd *cobra.Command, args []string) error {
	if len(args) < 1 {
		return errors.New("not enough arguments: command requires an argument")
	}
	opts.srpm = args[0]
	return nil
}

func (opts Options) Run() {
	args := make([]string, 0)
	if len(opts.buildroot) > 0 {
		args = append(args, "--root", opts.buildroot)
	}
	/*
		args = append(args, "--rebuild", opts.srpm)
		mock := rpak.MockCommand(context.TODO(), args...)
		mock.Stderr = os.Stderr
		mock.Stdout = os.Stdout
		if err := mock.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "Mock command failed: %s\n", err.Error())
			return
		}
	*/
}
