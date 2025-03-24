package linuxexec

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
)

func RunSudoCommand(command string, arguments ...string) error {
	return RunCommand("sudo", append([]string{command}, arguments...)...)
}

func RunCommand(command string, arguments ...string) error {
	var stdout, stderr bytes.Buffer

	cmd := exec.Command(command, arguments...)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	cmd.Stdin = os.Stdin

	if err := cmd.Run(); err != nil {
		// Collect output
		outStr := stdout.String()
		errStr := stderr.String()
		return fmt.Errorf("Error while running command '%s %v'. Error: %w\nStdout: %s\nStderr: %s", command, arguments, err, outStr, errStr)
	}

	return nil
}
