package main

import (
	"errors"
	"os"
	"os/exec"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	cmdExec := cmd[0]
	command := exec.Command(cmdExec, cmd[1:]...)

	for name, value := range env {
		_, ok := os.LookupEnv(name)

		if ok {
			err := os.Unsetenv(name)
			if err != nil {
				return 1
			}
		}
		if !value.NeedRemove {
			err := os.Setenv(name, value.Value)
			if err != nil {
				return 1
			}
		}
	}
	command.Env = os.Environ()
	command.Stdout = os.Stdout
	command.Stdin = os.Stdin

	err := command.Start()
	if err != nil {
		return 1
	}

	if err = command.Wait(); err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			return exitErr.ExitCode()
		}
		return 1
	}

	return 0
}
