package main

import (
	"os"
	"os/exec"
)

const OkCode = 0

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	name, args := cmd[0], cmd[1:]

	exCmd := exec.Command(name, args...)

	SetEnv(env)

	exCmd.Stdin, exCmd.Stdout, exCmd.Stderr = os.Stdin, os.Stdout, os.Stderr
	if err := exCmd.Run(); err != nil {
		var errCode *exec.ExitError
		return errCode.ExitCode()
	}

	return OkCode
}
