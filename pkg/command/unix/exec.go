package unix

import (
	"fmt"
	"log/slog"
	"os"
	"os/exec"
)

// runSudoCommand runs a linux command with sudo
func RunSudoCommand(command string, arguments ...string) error {
	return RunCommand("sudo", append([]string{command}, arguments...)...)
}

// runCommand will run a command in the shell
func RunCommand(command string, arguments ...string) error {
	return RunCommandWithStdin(os.Stdin, command, arguments...)
}

// RunCommandWithStdin will run a command in the shell and allow you to pass stdin
// @NOTE: If `stdin` is nil, then it will not be used. Hence, you can use this to run commands interactively or not
func RunCommandWithStdin(stdin *os.File, command string, arguments ...string) error {
	cmd := exec.Command(command, arguments...)
	if stdin != nil {
		cmd.Stdin = stdin
	}

	var out []byte
	var err error

	if out, err = cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("error while running command '%s %v'. Error: %w\nStdout: %s", command, arguments, err, out)
	}

	slog.Debug(string(out))

	return nil
}
