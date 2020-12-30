package initialize

import (
	"github.com/spf13/cobra"
	"encoding/json"
	"os"
	"code.pikelabs.net/go/pkg/soda"
	"strings"
)

type Options struct {
	specFile string
	sources string
}


func NewInitCmd() *cobra.Command {
	var opts Options
	c := &cobra.Command{
		Use: "init",
		Short: "initialize",
		Run: func(cmd *cobra.Command, args []string){
			opts.Prep()
			opts.Run()
		},
	}

	c.Flags().StringVarP(&opts.specFile, "spec", "s", "", "Specfile")
	c.Flags().StringVarP(&opts.sources, "sources", "d", "", "Sources")
	return c
}

func (opts *Options) Prep() error {
	return nil
}

func (opts *Options) Run() error {
	pkg := soda.Package{
		Specfile: opts.specFile,
		Sources: strings.Split(opts.sources, ","),
	}

	f, err := os.Create(".sodafile")
	if err != nil {
		return err
	}
	defer f.Close()
	codec := json.NewEncoder(f)
	codec.SetIndent("", "  ")
	err = codec.Encode(&pkg)
	return err
}