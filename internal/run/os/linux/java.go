package linux

import (
	"fmt"
	"os/exec"

	"github.com/Michaelpalacce/go-btva/internal/options"
	"github.com/Michaelpalacce/go-btva/internal/run/os/software"
	"github.com/Michaelpalacce/go-btva/pkg/command/unix"
	"github.com/Michaelpalacce/go-btva/pkg/os"
)

// JavaSoftware is responsible for installing and checking if java is installed
type JavaSoftware struct {
	os      *os.OS
	options *options.RunOptions

	initialized bool
}

var javaSoftware *JavaSoftware = &JavaSoftware{}

// Install will install java with apt
func (s *JavaSoftware) Install() error {
	return unix.RunSudoCommand("apt", "install", "-y", fmt.Sprintf("openjdk-%s-jdk", s.options.Software.JavaVersion))
}

// Exists verifies if java is already installed.
// `java --version` will return 0 if java exists and is setup correctly
// For example it will return `1` in case that `JAVA_HOME` is not configured ok
func (s *JavaSoftware) Exists() bool {
	cmd := exec.Command("java", "--version")
	_, err := cmd.CombinedOutput()

	return err == nil
}

func (s *JavaSoftware) GetName() string    { return software.JavaSoftwareKey }
func (s *JavaSoftware) GetVersion() string { return s.options.Software.JavaVersion }

// Java will return the JavaSoftware object that can be used to install and check if java exists
// Only a single instance of the JavaSoftware will be returned
func (i *Installer) Java() software.Software {
	if !javaSoftware.initialized {
		javaSoftware.os = i.OS
		javaSoftware.options = i.Options
		javaSoftware.initialized = true
	}

	return javaSoftware
}
