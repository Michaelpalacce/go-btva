package linux

import (
	"fmt"
	osz "os"
	"os/exec"

	"github.com/Michaelpalacce/go-btva/internal/args"
	"github.com/Michaelpalacce/go-btva/internal/software"
	"github.com/Michaelpalacce/go-btva/pkg/os"
)

// NodeSoftware is responsible for installing, removing and checking if node is installed
type NodeSoftware struct {
	OS      *os.OS
	Options *args.Options

	initialized bool
}

var nodeSoftware *NodeSoftware = &NodeSoftware{}

// Install will install node with apt
func (s *NodeSoftware) Install() error {
	shell := osz.Getenv("SHELL")
	if shell == "" {
		shell = "/bin/bash"
	}

	if err := RunCommand(shell, "-c", fmt.Sprintf("curl -fsSL https://fnm.vercel.app/install | %s", shell)); err != nil {
		return err
	}

	var profile string
	switch shell {
	case "/bin/zsh":
		profile = "$HOME/.zshrc"
	case "/bin/fish":
		profile = "$HOME/.config/fish/config.fish"
	default:
		profile = "$HOME/.bashrc"
	}

	return RunCommand(shell, "-i", "-c", fmt.Sprintf("source %s && fnm install %s", profile, s.Options.Software.LinuxNodeVersion))
}

// Exists verifies if node is already installed.
// Relies on `which`, which returns exit code 0 if the program is found and 1 if not
func (s *NodeSoftware) Exists() bool {
	_, err := exec.Command("which", "node").CombinedOutput()

	return err == nil
}

func (s *NodeSoftware) GetName() string    { return software.NodeSoftwareKey }
func (s *NodeSoftware) GetVersion() string { return s.Options.Software.LinuxNodeVersion }

// Node will return the NodeSoftware object that can be used to install, remove or check if node exists
// Only a single instance of the NodeSoftware will be returned
func (i *LinuxInstaller) Node() software.Software {
	if !nodeSoftware.initialized {
		nodeSoftware.OS = i.OS
		nodeSoftware.Options = i.Options
		nodeSoftware.initialized = true
	}

	return nodeSoftware
}
