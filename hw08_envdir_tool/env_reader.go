package main

import (
	"bufio"
	"bytes"
	"os"
	"path/filepath"
	"strings"
)

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	dirEntries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	env := make(Environment, len(dirEntries))

	for _, dirE := range dirEntries {
		envVarName := dirE.Name()
		if strings.Contains(envVarName, "=") {
			continue
		}

		fileInfo, err := dirE.Info()
		if err != nil {
			return nil, err
		}

		isNeedToRemove := fileInfo.Size() == 0

		file, err := os.Open(filepath.Join(dir, dirE.Name()))
		if err != nil {
			return nil, err
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)

		scanner.Scan()
		buf := scanner.Bytes()

		if err != nil {
			return nil, err
		}

		buf = bytes.Split(buf, []byte{0x0A})[0]
		buf = bytes.ReplaceAll(buf, []byte{0x00}, []byte{0x0A})

		envVarValue := string(buf)
		envVarValue = strings.TrimRight(envVarValue, " \t")

		env[envVarName] = EnvValue{
			envVarValue,
			isNeedToRemove,
		}
	}

	return env, nil
}
