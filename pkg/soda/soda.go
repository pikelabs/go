package soda

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
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
