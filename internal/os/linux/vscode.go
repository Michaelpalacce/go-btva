package linux

import (
	"os/exec"

	"github.com/Michaelpalacce/go-btva/internal/args"
	"github.com/Michaelpalacce/go-btva/internal/os/software"
	"github.com/Michaelpalacce/go-btva/pkg/command/unix"
	"github.com/Michaelpalacce/go-btva/pkg/os"
)

// VsCodeSoftware is responsible for installing, removing and checking if VSCode is installed
type VsCodeSoftware struct {
	os      *os.OS
	options *args.Options

	initialized bool
}

var vsCodeSoftware *VsCodeSoftware = &VsCodeSoftware{}

// Install will install vsCode with apt
func (s *VsCodeSoftware) Install() error {
	if err := unix.RunSudoCommand("snap", "install", "--classic", "code"); err != nil {
		return err
	}
	extensions := []string{
		"vmware-pscoe.vrealize-developer-tools",
	}

	for _, extension := range extensions {
		if err := unix.RunCommand("code", "--install-extension", extension); err != nil {
			return err
		}
	}

	return nil
}

// Exists verifies if vsCode is already installed.
// `vsCode --version` will return 0 if vsCode exists and is setup correctly
// For example it will return `1` in case that `JAVA_HOME` is not configured ok
func (s *VsCodeSoftware) Exists() bool {
	cmd := exec.Command("which", "code")
	_, err := cmd.CombinedOutput()

	return err == nil
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
