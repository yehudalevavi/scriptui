package script

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func LoadScripts() (map[string]*Script, error) {
	exePath, err := os.Executable()
	if err != nil {
		return nil, fmt.Errorf("failed to get executable path: %v\n", err)
	}
	baseDir := filepath.Dir(exePath)

	scriptsDir := filepath.Join(baseDir, "scripts")
	argsDir := filepath.Join(baseDir, "args")

	scriptFiles, err := os.ReadDir(scriptsDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read scripts directory: %v\n", err)
	}

	ret := make(map[string]*Script, 10)
	for _, scriptFile := range scriptFiles {
		if scriptFile.IsDir() {
			continue
		}

		scriptPath := filepath.Join(scriptsDir, scriptFile.Name())
		scriptBaseName := strings.TrimSuffix(scriptFile.Name(), filepath.Ext(scriptFile.Name()))
		argFilePath := filepath.Join(argsDir, scriptBaseName+".json")

		s, err := New(scriptPath, argFilePath, scriptBaseName)
		if err != nil {
			return nil, fmt.Errorf("failed to create script: %v\n", err)
		}
		ret[s.Name] = s
	}
	return ret, nil
}
