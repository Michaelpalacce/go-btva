package linux

import (
	"fmt"
	"os/exec"

	"github.com/Michaelpalacce/go-btva/internal/options"
	"github.com/Michaelpalacce/go-btva/internal/run/os/software"
	"github.com/Michaelpalacce/go-btva/pkg/command/unix"
	"github.com/Michaelpalacce/go-btva/pkg/os"
)

// NodeSoftware is responsible for installing and checking if node is installed
type NodeSoftware struct {
	os      *os.OS
	options *options.RunOptions

	initialized bool
}

var nodeSoftware *NodeSoftware = &NodeSoftware{}

// Install will install node with the help of `fnm`
// After the execution of the tool, you will need to either source your zsh file or open a new terminal
func (s *NodeSoftware) Install() error {
	shell := s.os.Shell
	profile := s.os.ShellProfile

	if err := unix.RunCommand(shell, "-c", fmt.Sprintf("curl -fsSL https://fnm.vercel.app/install | %s", shell)); err != nil {
		return err
	}

	return unix.RunCommand(shell, "-i", "-c", fmt.Sprintf("source %s && fnm install %s", profile, s.options.Software.NodeVersion))
}

// Exists verifies if node is already installed.
// Relies on `which`, which returns exit code 0 if the program is found and 1 if not
func (s *NodeSoftware) Exists() bool {
	_, err := exec.Command("which", "node").CombinedOutput()

	return err == nil
}

func (s *NodeSoftware) GetName() string    { return software.NodeSoftwareKey }
func (s *NodeSoftware) GetVersion() string { return s.options.Software.NodeVersion }

// Node will return the NodeSoftware object that can be used to install and check if node exists
// Only a single instance of the NodeSoftware will be returned
func (i *Installer) Node() software.Software {
	if !nodeSoftware.initialized {
		nodeSoftware.os = i.OS
		nodeSoftware.options = i.Options
		nodeSoftware.initialized = true
	}

	return nodeSoftware
}
