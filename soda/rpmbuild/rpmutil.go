package rpmbuild

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
)

const (
	SpecDir    = "SPECS"
	SourcesDir = "SOURCES"
	SRPMDir    = "SRPMS"
	RPMSDir    = "RPMS"
)

const defaultDirMode = 0700

func mkdir(path string, fileMode os.FileMode) error {
	stats, err := os.Stat(path)
	if err != nil {
		if !os.IsNotExist(err) {
			return err
		}
		if err = os.Mkdir(path, fileMode); err != nil {
			return err
		}
		if stats, err = os.Stat(path); err != nil {
			return err
		}
	}
	mode := stats.Mode()
	if !mode.IsDir() {
		return errors.New("can't create directiry %s: File exists")
	}
	return nil
}

// copy regular  file to desitnation (directory or file)
func CopyFile(dest, src string) error {
	fmt.Printf("%s -> %s\n", src, dest)
	destStats, err := os.Stat(dest)
	if err != nil {
		return err
	}
	finalDest := dest
	if destStats.IsDir() {
		fName := filepath.Base(src)
		finalDest = path.Join(dest, fName)
	}

	src, err = filepath.Abs(src)
	if err != nil {
		return err
	}

	srcStats, err := os.Stat(src)
	if err != nil {
		return err
	}

	switch mode := srcStats.Mode(); {
		case mode.IsRegular():
			_, err = copyFile(finalDest, src)
			return err
		case mode.IsDir():
			list, err := ioutil.ReadDir(src)
			if err != nil {
				return err
			}
			for _, info := range list {
				fmt.Printf("processing %s\n", info.Name())
				if strings.HasPrefix(info.Name(), ".") {
					continue
				}
				p := path.Join(src, info.Name())
				fmt.Printf("path %s\n", p)
				mode := info.Mode()
				if mode.IsRegular() {
					err = CopyFile(finalDest, p)
				} else {
					err = fmt.Errorf("can't copy nested dirstories")
				}
				if err != nil {
					break
				}
			}
			return err

		default:
			return fmt.Errorf("supporting only regular files or directories")
		}
}

func copyFile(dest, src string) (int64, error) {
	srcStats, err := os.Stat(src)
	if err != nil {
		return 0, err
	}
	if !srcStats.Mode().IsRegular() {
		return 0, fmt.Errorf("%s isn't regular file", src)
	}
	in, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer in.Close()

	out, err := os.Create(dest)
	if err != nil {
		return 0, err
	}

	defer out.Close()


	n, err := io.Copy(out, in)
	if err != nil {
		return n, err
	}


	err = out.Sync()
	if err != nil {
		return n, err
	}
	err = os.Chmod(dest, srcStats.Mode())
	
	return n, err
}

// MkBuildroot creates directory layout for rpmbuild
func MkBuildroot(workdir, spec string, sources []string) error {
	specDir := path.Join(workdir, SpecDir)
	srcDir := path.Join(workdir, SourcesDir)
	if err := mkdir(specDir, defaultDirMode); err != nil {
		return err
	}
	if err := mkdir(srcDir, defaultDirMode); err != nil {
		return err
	}
	if err := CopyFile(specDir, spec); err != nil {
		return err
	}

	for _, f := range sources {
		if err := CopyFile(srcDir, f); err != nil {
			return err
		}
	}

	return nil
}

