package windows

import (
	"github.com/Michaelpalacce/go-btva/internal/options"
	"github.com/Michaelpalacce/go-btva/internal/run/os/software"
	"github.com/Michaelpalacce/go-btva/pkg/os"
)

// JavaSoftware is responsible for installing and checking if java is installed
type JavaSoftware struct {
	os      *os.OS
	options *options.RunOptions

	initialized bool
}

var javaSoftware *JavaSoftware = &JavaSoftware{}

func (s *JavaSoftware) Install() error {
	return nil
}

func (s *JavaSoftware) Exists() bool {
	return true
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
