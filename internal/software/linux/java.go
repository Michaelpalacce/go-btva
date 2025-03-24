package linux

import (
	"fmt"
	"os/exec"

	"github.com/Michaelpalacce/go-btva/internal/args"
	"github.com/Michaelpalacce/go-btva/internal/software"
	"github.com/Michaelpalacce/go-btva/pkg/os"
)

// JavaSoftware is responsible for installing, removing and checking if java is installed
type JavaSoftware struct {
	os      *os.OS
	options *args.Options

	initialized bool
}

var javaSoftware *JavaSoftware = &JavaSoftware{}

// Install will install java with apt
func (s *JavaSoftware) Install() error {
	return RunSudoCommand("apt", "install", "-y", fmt.Sprintf("openjdk-%s-jdk", s.options.Software.LinuxJavaVersion))
}

// Exists verifies if java is already installed.
// Relies on `which`, which returns exit code 0 if the program is found and 1 if not
func (s *JavaSoftware) Exists() bool {
	cmd := exec.Command("which", "java")
	_, err := cmd.CombinedOutput()

	return err == nil
}

func (s *JavaSoftware) GetName() string    { return software.JavaSoftwareKey }
func (s *JavaSoftware) GetVersion() string { return s.options.Software.LinuxJavaVersion }

// Java will return the JavaSoftware object that can be used to install, remove or check if java exists
// Only a single instance of the JavaSoftware will be returned
func (i *LinuxInstaller) Java() software.Software {
	if !javaSoftware.initialized {
		javaSoftware.os = i.OS
		javaSoftware.options = i.Options
		javaSoftware.initialized = true
	}

	return javaSoftware
}
