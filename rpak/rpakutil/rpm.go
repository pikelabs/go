package rpakutil

import (
	"errors"
	"io"
	"os"
	"os/exec"
	"path/filepath"
)

var (
	ErrNotRegularFile = errors.New("Not a regular file")
)

// ExtractSRPM exports contents of SRPM to directory
//
func ExtractSRPM(srpm string, dir string) error {
	srpm, err := filepath.Abs(srpm)
	if err != nil {
		return err
	}
	finfo, err := os.Stat(srpm)
	if err != nil {
		return err
	}
	if !finfo.Mode().IsRegular() {
		return ErrNotRegularFile
	}

	// XXX: for now using external processes to extract
	// rpm contents, later on we should switch to internal
	// implementation
	rp, wp := io.Pipe()
	cmd := exec.Command("rpm2cpio", srpm)
	cmd.Stdout = wp
	cmd2 := exec.Command("cpio", "-iud", "--quiet")
	cmd2.Dir = dir
	cmd2.Stdin = rp
	go func() {
		err := cmd.Run()
		wp.CloseWithError(err)
	}()

	return cmd2.Run()
}
