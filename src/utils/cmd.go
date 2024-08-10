package utils

import (
	"fmt"
	"os/exec"
	"strings"
)

func ExecuteCommand(command, dir string) (string, error) {

	parts := strings.Fields(command)

	if len(parts) == 0 {
		return "", fmt.Errorf("no command provided")
	}

	cmd := exec.Command(parts[0], parts[1:]...)

	if dir != "" {
		cmd.Dir = dir
	}

	output, err := cmd.CombinedOutput() //  stdout & stderr
	if err != nil {
		return string(output), err
	}

	return string(output), nil
}
