package soda

import (
	"archive/tar"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"strings"
)

type Package struct {
	Specfile      string                  `json:"specfile"`
	VersionCmd    string                  `json:"version_cmd"`
	PrebuildSteps []ExecutableStepEnvelop `json:"prebuild"`
	Sources       []string                `json:"sources"`

	version []byte `json:-`
}

type ExecutableStep interface {
	Run() error
}

type ExecutableStepEnvelop struct {
	Type string          `json: "type"`
	Msg  json.RawMessage `json:"options"`
}

func (s ExecutableStepEnvelop) Run() error {
	var realStep ExecutableStep
	switch s.Type {
	case "shell":
		realStep = &ShellStep{}
	case "tarball":
		realStep = &TarballStep{}
	default:
		return fmt.Errorf("unknown type: %s", s.Type)
	}
	if err := json.Unmarshal(s.Msg, realStep); err != nil {
		return err
	}
	return realStep.Run()
}

type ShellStep struct {
	Cmd string `json:"cmd"`
}

func (s *ShellStep) Run() error {
	return Shellout(s.Cmd)
}

type TarballStep struct {
	Name     string   `json:"name"`
	List     []string `json:"files"`
	Basepath string   `json:"base_path"`
	Excludes []string `json:"excludes"`
}

func (s *TarballStep) Run() error {
	f, err := os.Create(s.Name)
	if err != nil {
		return err
	}
	defer f.Close()

	fileList := make([]string, len(s.List))

	for _, f := range s.List {
		matched, err := matchedAny(f, s.Excludes)
		if err != nil {
			return err
		}
		if matched {
			continue
		}
		info, err := os.Stat(f)
		if err != nil {
			return err
		}
		if info.IsDir() {
			err = filepath.Walk(f, func(p string, info os.FileInfo, err error) error {
				if info.IsDir() && strings.HasPrefix(info.Name(), ".") {
					return filepath.SkipDir
				}

				if strings.HasPrefix(info.Name(), ".") {
					return nil
				}

				fileList = append(fileList, p)
				return nil

			})
			if err != nil {
				return err
			}
		} else {
			fileList = append(fileList, f)
		}
	}

	return createTarball(fileList, s.Basepath, f)
}

func matchedAny(s string, patterns []string) (bool, error) {
	for _, pattern := range patterns {
		matched, err := regexp.MatchString(pattern, s)
		if err != nil {
			return false, err
		}
		if matched {
			return true, nil
		}
	}
	return false, nil
}

func createTarball(files []string, basepath string, w io.Writer) error {
	tw := tar.NewWriter(w)
	defer tw.Close()
	for _, f := range files {
		if err := addToTar(tw, f, basepath); err != nil {
			return err
		}
	}
	return nil
}

func addToTar(w *tar.Writer, filename, basepath string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		return err
	}

	header, err := tar.FileInfoHeader(info, info.Name())
	if err != nil {
		return err
	}
	header.Name = path.Join(basepath, info.Name())

	if err = w.WriteHeader(header); err != nil {
		return err
	}
	_, err = io.Copy(w, file)

	return err
}

func (p *Package) GetVersion() (v []byte, err error) {
	if p.version != nil {
		return p.version, nil
	}
	v, err = exec.Command(p.VersionCmd).Output()
	if err != nil {
		return nil, err
	}
	p.version = v
	return
}

func Shellout(cmd string) error {
	return RunCmd("/bin/sh", "-c", cmd)
}

func RunCmd(name string, args ...string) error {
	fmt.Printf("running %s %s\n", name, strings.Join(args, " "))
	c := exec.Command(name, args...)
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	return c.Run()
}
