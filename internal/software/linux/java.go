package linux

import (
	"fmt"
	"log/slog"
	"os/exec"

	"github.com/Michaelpalacce/go-btva/internal/args"
	"github.com/Michaelpalacce/go-btva/internal/software"
	"github.com/Michaelpalacce/go-btva/pkg/os"
)

// JavaSoftware is responsible for installing, removing and checking if java is installed
type JavaSoftware struct {
	OS      *os.OS
	Options *args.Options

	initialized bool
}

var javaSoftware *JavaSoftware = &JavaSoftware{}

// Install will install java with apt
// @NOTE: Make sure you run the go process as sudo
func (s *JavaSoftware) Install() error {
	cmd := exec.Command("apt", "install", "-y", fmt.Sprintf("openjdk-%s-jdk", s.Options.Software.LinuxJavaVersion))

	output, err := cmd.CombinedOutput()
	slog.Debug(string(output))

	return err
}

// Remove will remove java
func (s *JavaSoftware) Remove() error {
	cmd := exec.Command("apt", "remove", "-y", fmt.Sprintf("openjdk-%s-jdk", s.Options.Software.LinuxJavaVersion))

	output, err := cmd.CombinedOutput()
	slog.Debug(string(output))

	return err
}

// Exists verifies if java is already installed.
// Relies on `which`, which returns exit code 0 if the program is found and 1 if not
func (s *JavaSoftware) Exists() bool {
	cmd := exec.Command("which", "java")
	_, err := cmd.CombinedOutput()

	return err == nil
}

// Java will return the JavaSoftware object that can be used to install, remove or check if java exists
// Only a single instance of the JavaSoftware will be returned
func (i *LinuxInstaller) Java() software.Software {
	if !javaSoftware.initialized {
		javaSoftware.OS = i.OS
		javaSoftware.Options = i.Options
		javaSoftware.initialized = true
	}

	return javaSoftware
}
