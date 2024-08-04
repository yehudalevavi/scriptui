package script

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
)

type Script struct {
	Path    string `json:"path"`
	Name    string `json:"name"`
	Version string `json:"version"`
	Args    []Arg  `json:"args"`
	Desc    string `json:"desc"`
}

type Arg struct {
	Flag         string      `json:"flag"`
	Value        interface{} `json:"value"`
	Mandatory    bool        `json:"mandatory"`
	HasValue     bool        `json:"hasValue"`
	Title        string      `json:"title"`
	DefaultValue string      `json:"defaultValue"`
	Placeholder  string      `json:"placeholder"`
}

func New(scriptPath, argsFilePath, scriptBaseName string) (*Script, error) {
	jsonData, err := os.ReadFile(argsFilePath)
	if err != nil {
		return nil, err
	}

	ret := &Script{}
	if err := json.Unmarshal(jsonData, &ret); err != nil {
		return nil, err
	}

	ret.Path = scriptPath
	ret.Name = scriptBaseName

	return ret, nil
}

func (s *Script) CreateCmd() (*exec.Cmd, error) {
	args, err := s.PrepareArgs()
	fmt.Println(args)
	if err != nil {
		return nil, err
	}
	cmd := exec.Command("bash", args...)
	return cmd, nil
}

func (s *Script) PrepareArgs() ([]string, error) {
	args := []string{s.Path}
	if len(args) == 0 {
		return args, nil
	}
	for _, arg := range s.Args {
		if arg.Mandatory && arg.Value == nil {
			return nil, fmt.Errorf("mandatory arg %s cannot be nil", arg.Flag)
		}
		if arg.Value != nil {
			if arg.Flag != "" {
				args = append(args, arg.Flag, fmt.Sprintf("%v", arg.Value))
			} else {
				args = append(args, fmt.Sprintf("%v", arg.Value))
			}
		} else if arg.Value == nil && !arg.Mandatory {
			args = append(args, arg.Flag)
		}
	}
	return args, nil
}

func (s *Script) RunScript() {
	// Reset the terminal using the 'reset' command
	cmd := exec.Command("reset")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	err := cmd.Run()
	if err != nil {
		// Fallback to 'stty sane' if 'reset' fails
		cmd = exec.Command("stty", "sane")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin
		_ = cmd.Run()
	}
	cmd, err = s.CreateCmd()
	if err != nil {
		panic(err)
	}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Run()
	os.Exit(0)
}
