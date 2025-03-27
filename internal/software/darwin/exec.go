package darwin

import (
	"fmt"
	"log/slog"
	"os"
	"os/exec"
)

// runSudoCommand runs a darwin command with sudo
func runSudoCommand(command string, arguments ...string) error {
	return runCommand("sudo", append([]string{command}, arguments...)...)
}

// runCommand will run a command in the shell
func runCommand(command string, arguments ...string) error {
	cmd := exec.Command(command, arguments...)
	cmd.Stdin = os.Stdin

	var out []byte
	var err error

	if out, err = cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("error while running command '%s %v'. Error: %w\nStdout: %s", command, arguments, err, out)
	}

	slog.Debug(string(out))

	return nil
}
