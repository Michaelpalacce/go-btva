package darwin

import (
	"fmt"
	osz "os"
	"os/exec"

	"github.com/Michaelpalacce/go-btva/internal/args"
	"github.com/Michaelpalacce/go-btva/internal/software"
	"github.com/Michaelpalacce/go-btva/pkg/command/unix"
	"github.com/Michaelpalacce/go-btva/pkg/os"
)

// NodeSoftware is responsible for installing, removing and checking if node is installed
type NodeSoftware struct {
	os      *os.OS
	options *args.Options

	initialized bool
}

var nodeSoftware *NodeSoftware = &NodeSoftware{}

// Install will install node with apt
func (s *NodeSoftware) Install() error {
	shell := osz.Getenv("SHELL")
	if shell == "" {
		shell = "/bin/bash"
	}

	var profile string
	switch shell {
	case "/bin/zsh":
		profile = "$HOME/.zshrc"
	case "/bin/fish":
		profile = "$HOME/.config/fish/config.fish"
	case "/bin/bash":
		profile = "$HOME/.bashrc"
	default:
		return fmt.Errorf("Shell %s is not supported", shell)
	}

	if err := unix.RunCommand("brew", "install", "fnm"); err != nil {
		return err
	}

	if err := unix.RunCommand(shell, "-c", fmt.Sprintf("echo 'eval \"$(fnm env --use-on-cd)\"' >> %s", profile)); err != nil {
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

// Node will return the NodeSoftware object that can be used to install, remove or check if node exists
// Only a single instance of the NodeSoftware will be returned
func (i *DarwinInstaller) Node() software.Software {
	if !nodeSoftware.initialized {
		nodeSoftware.os = i.OS
		nodeSoftware.options = i.Options
		nodeSoftware.initialized = true
	}

	return nodeSoftware
}
