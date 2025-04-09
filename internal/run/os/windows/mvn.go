package windows

import (
	"github.com/Michaelpalacce/go-btva/internal/options"
	"github.com/Michaelpalacce/go-btva/internal/run/os/software"
	"github.com/Michaelpalacce/go-btva/pkg/os"
)

// MvnSoftware is responsible for installing and checking if mvn is installed
type MvnSoftware struct {
	os      *os.OS
	options *options.RunOptions

	initialized bool
}

var mvnSoftware *MvnSoftware = &MvnSoftware{}

func (s *MvnSoftware) Install() error {
	return nil
}

func (s *MvnSoftware) Exists() bool {
	return true
}

func (s *MvnSoftware) GetName() string    { return software.MvnSoftwareKey }
func (s *MvnSoftware) GetVersion() string { return s.options.Software.MvnVersion }

// Java will return the MvnSoftware object that can be used to install and check if mvn exists
// Only a single instance of the MvnSoftware will be returned
func (i *Installer) Mvn() software.Software {
	if !mvnSoftware.initialized {
		mvnSoftware.os = i.OS
		mvnSoftware.options = i.Options
		mvnSoftware.initialized = true
	}

	return mvnSoftware
}
