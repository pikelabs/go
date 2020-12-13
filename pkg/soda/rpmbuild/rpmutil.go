package rpmbuild

import (
	"errors"
	"os"
	"path"
)

const (
	SpecDir    = "SPEC"
	SourcesDir = "SOURCES"
	SRPMDir    = "SRPMS"
	RPMSDir    = "RPMS"
)

const defaultFileMode = 0600

func mkdir(path string, fileMode os.FileMode) error {
	stats, err := os.Stat(path)
	if err != nil {
		if !os.IsNotExist(err) {
			return err
		}
		if err = os.Mkdir(path, fileMode); err != nil {
			return err
		}
	}
	mode := stats.Mode()
	if !mode.IsDir() {
		return errors.New("can't create directiry %s: File exists")
	}
	return nil
}

// MkBuildroot creates directory layout for rpmbuild
func MkBuildroot(workdir, spec string, sources []string) error {
	specDir := path.Join(workdir, SpecDir)
	srcDir := path.Join(workdir, SourcesDir)
	if err := mkdir(specDir, defaultFileMode); err != nil {
		return err
	}
	if err := mkdir(srcDir, defaultFileMode); err != nil {
		return err
	}
	return nil
}
