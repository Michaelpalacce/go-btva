package windows

import (
	"github.com/Michaelpalacce/go-btva/internal/options"
	"github.com/Michaelpalacce/go-btva/internal/run/os/software"
	"github.com/Michaelpalacce/go-btva/pkg/os"
)

// VsCodeSoftware is responsible for installing and checking if VSCode is installed
type VsCodeSoftware struct {
	os      *os.OS
	options *options.RunOptions

	initialized bool
}

var vsCodeSoftware *VsCodeSoftware = &VsCodeSoftware{}

func (s *VsCodeSoftware) Install() error {
	return nil
}

func (s *VsCodeSoftware) Exists() bool {
	return true
}

func (s *VsCodeSoftware) GetName() string    { return software.VsCodeSoftwareKey }
func (s *VsCodeSoftware) GetVersion() string { return s.options.Software.VsCodeVersion }

// VsCode will return the VsCodeSoftware object that can be used to install and check if vsCode exists
// Only a single instance of the VsCodeSoftware will be returned
func (i *Installer) VsCode() software.Software {
	if !vsCodeSoftware.initialized {
		vsCodeSoftware.os = i.OS
		vsCodeSoftware.options = i.Options
		vsCodeSoftware.initialized = true
	}

	return vsCodeSoftware
}
