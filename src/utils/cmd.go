package utils

import (
	"os"
	"os/exec"
)

func ExecuteCommand(command string, args, env []string, dir string) (string, error) {

	cmd := exec.Command(command, args...)

	if dir != "" {
		cmd.Dir = dir
	}

	cmd.Env = os.Environ()

	if len(env) > 0 {
		cmd.Env = append(cmd.Env, env...)
	}

	output, err := cmd.CombinedOutput() //  stdout & stderr
	if err != nil {
		return string(output), err
	}

	return string(output), nil
}
