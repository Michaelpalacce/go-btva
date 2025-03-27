package darwin

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
)

// runSudoCommand runs a darwin command with sudo
func runSudoCommand(command string, arguments ...string) error {
	return runCommand("sudo", append([]string{command}, arguments...)...)
}

// runCommand will run a command in the shell
func runCommand(command string, arguments ...string) error {
	var stdout, stderr bytes.Buffer

	cmd := exec.Command(command, arguments...)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	cmd.Stdin = os.Stdin

	if err := cmd.Run(); err != nil {
		// Collect output
		outStr := stdout.String()
		errStr := stderr.String()
		return fmt.Errorf("error while running command '%s %v'. Error: %w\nStdout: %s\nStderr: %s", command, arguments, err, outStr, errStr)
	}

	return nil
}
