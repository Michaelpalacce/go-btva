package linux

import (
	"log/slog"
	"os/exec"
)

func (i *LinuxInstaller) InstallJava() error {
	cmd := exec.Command("apt", "install", "-y", i.Options.Software.JavaLinuxPackage)

	output, err := cmd.CombinedOutput()
	slog.Debug(string(output))

	return err
}
