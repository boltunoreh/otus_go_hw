package main

import (
	"os"
)

func main() {
	argsWithProg := os.Args[1:]

	envDir := argsWithProg[0]
	commandWithArgs := argsWithProg[1:]

	env, err := ReadDir(envDir)
	if err != nil {
		return
	}

	RunCmd(commandWithArgs, env)
}
