package darwin

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

// Install will install java with brew. Brew will also configure the JAVA_HOME for us
func (s *JavaSoftware) Install() error {
	return unix.RunCommand("brew", "install", fmt.Sprintf("openjdk@%s", s.options.Software.JavaVersion))
}

// Exists verifies if java is already installed.
// Relies on `which`, which returns exit code 0 if the program is found and 1 if not
func (s *JavaSoftware) Exists() bool {
	cmd := exec.Command("java", "-version")
	_, err := cmd.CombinedOutput()

	return err == nil
}

func (s *JavaSoftware) GetName() string    { return software.JavaSoftwareKey }
func (s *JavaSoftware) GetVersion() string { return s.options.Software.JavaVersion }

// Java will return the JavaSoftware object that can be used to install, remove or check if java exists
// Only a single instance of the JavaSoftware will be returned
func (i *Installer) Java() software.Software {
	if !javaSoftware.initialized {
		javaSoftware.os = i.OS
		javaSoftware.options = i.Options
		javaSoftware.initialized = true
	}

	return javaSoftware
}
