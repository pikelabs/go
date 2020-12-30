package rpak

import (
	"errors"
	"os"
	"path/filepath"

	"code.pikelabs.net/go/rpak/rpakutil"
	"code.pikelabs.net/go/rpm/rpmutil"
)

var (
	ErrUncommitedChanges = errors.New("Repo has uncommited changes")
	ErrNotRegularFile    = errors.New("Not a regular file")
)

type Layout struct {
	Rootdir   string
	Specdir   string
	Sourcedir string
}

var defaultLayout = Layout{
	Rootdir:   ".",
	Specdir:   "SPECS",
	Sourcedir: "SOURCES",
}

type Workflow struct {
	Path      string
	Repo      Repo
	Layout    Layout
	Lookaside Lookaside
}

func NewWorkflow(path string) *Workflow {
	return &Workflow{
		Path:      path,
		Layout:    defaultLayout,
		Repo:      &GitRepo{},
		Lookaside: &GitLFS{},
	}
}

func CloneRepo(url string) *Workflow {
	return nil
}

// LoadRepo loads repo information from path
func (w *Workflow) LoadRepo() (err error) {
	return w.Repo.Open(w.Path)
}

// ImportSRPM loads SRPM to the path
func (w Workflow) ImportSRPM(path string) error {
	isClean, err := w.Repo.IsClean()
	if err != nil {
		return err
	}
	if !isClean {
		return ErrUncommitedChanges
	}

	path, err = filepath.Abs(path)
	if err != nil {
		return err
	}

	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()
	finfo, err := f.Stat()
	if err != nil {
		return err
	}
	if !finfo.Mode().IsRegular() {
		return ErrNotRegularFile
	}
	pkg, err := rpmutil.ReadPackage(f)

	if err != nil {
		return err
	}
	contents, err := pkg.Files()
	if err != nil {
		return err
	}
	cwd, _ := os.Getwd()
	rpakutil.ExtractSRPM(path, cwd)
	files := make([]string, 0, len(contents))
	for _, f := range contents {
		if ok, _ := w.Lookaside.Eligable(f.Name); !ok {
			files = append(files, f.Name)
		} else {
			w.Lookaside.Upload(f.Name)
		}
	}

	err = w.Repo.Stage(files...)
	if err != nil {
		return err
	}
	return nil
}
