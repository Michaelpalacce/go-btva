package windows

import (
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"strings"
)

// RunElevatedCommand runs a Windows command with elevated privileges (triggers UAC prompt).
// It uses PowerShell's Start-Process -Verb RunAs.
// NOTE: Standard input redirection is generally NOT feasible or reliable for elevated
// commands launched this way. Therefore, this function does not support passing stdin.
func RunElevatedCommand(command string, arguments ...string) error {
	psCommandPath := strings.ReplaceAll(command, "'", "''")

	var psArgs []string
	for _, arg := range arguments {
		escapedArg := strings.ReplaceAll(arg, "'", "''")
		psArgs = append(psArgs, fmt.Sprintf("'%s'", escapedArg)) // Enclose in single quotes
	}
	psArgumentList := strings.Join(psArgs, ", ")

	psCommand := fmt.Sprintf("Start-Process -FilePath '%s' -ArgumentList @(%s) -Verb RunAs -Wait", psCommandPath, psArgumentList)

	slog.Debug("Attempting to run elevated command via PowerShell", "powershell_command", psCommand)

	err := RunCommandWithStdin(nil, "powershell", "-NoProfile", "-Command", psCommand)
	if err != nil {
		return fmt.Errorf("error attempting to elevate command '%s %v': %w", command, arguments, err)
	}
	return nil
}

// RunCommand will run a command in the system's shell (cmd.exe or PowerShell context).
func RunCommand(command string, arguments ...string) error {
	return RunCommandWithStdin(os.Stdin, command, arguments...)
}

// RunCommandWithStdin will run a command in the shell and allow you to pass stdin.
// @NOTE: If `stdin` is nil, then standard input will not be redirected to the command.
func RunCommandWithStdin(stdin *os.File, command string, arguments ...string) error {
	cmd := exec.Command(command, arguments...)
	if stdin != nil {
		cmd.Stdin = stdin
	}

	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("error running command '%s %v': %w\nOutput: %s", command, arguments, err, string(out))
	}

	outputStr := string(out)
	if len(outputStr) > 0 {
		slog.Debug("Command executed successfully", "command", command, "args", arguments, "output", outputStr)
	} else {
		slog.Debug("Command executed successfully (no output)", "command", command, "args", arguments)
	}

	return nil
}
