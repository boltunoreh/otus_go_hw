package main

import (
	"os"
	"os/exec"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	commandName := cmd[0]
	commandArgs := cmd[1:]
	command := exec.Command(commandName, commandArgs...)

	command.Stdin = os.Stdin
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr

	for envName, envValue := range env {
		var err error
		if envValue.NeedRemove {
			err = os.Unsetenv(envName)
		} else {
			err = os.Setenv(envName, envValue.Value)
		}

		if err != nil {
			return 1
		}
	}

	command.Env = os.Environ()

	if err := command.Run(); err != nil {
		return command.ProcessState.ExitCode()
	}

	return command.ProcessState.ExitCode()
}
