package rpak

import (
	"bytes"
	"os/exec"
	"strings"
)

type Lookaside interface {
	Init(path string) error
	Eligable(path string) (bool, error)
	Upload(path string) error
	Download(file string) error
}

type GitLFS struct {
}

func (lfs GitLFS) Init(path string) error {
	return nil
}

func (lfs GitLFS) Eligable(path string) (bool, error) {
	var output bytes.Buffer
	cmd := exec.Command("file", "-b", "--mime-encoding", path)
	cmd.Stdout = &output
	if err := cmd.Run(); err != nil {
		return false, err
	}
	mine := strings.TrimSpace(output.String())
	return mine == "binary", nil
}

func (lfs GitLFS) Upload(path string) error {
	cmd := exec.Command("git", "lfs", "track", path)
	if err := cmd.Run(); err != nil {
		return err
	}
	cmd = exec.Command("git", "add", path)
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}

func (lfs GitLFS) Download(path string) error {
	return nil
}
