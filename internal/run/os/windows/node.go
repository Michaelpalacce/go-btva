package windows

import (
	"github.com/Michaelpalacce/go-btva/internal/options"
	"github.com/Michaelpalacce/go-btva/internal/run/os/software"
	"github.com/Michaelpalacce/go-btva/pkg/os"
)

// NodeSoftware is responsible for installing and checking if node is installed
type NodeSoftware struct {
	os      *os.OS
	options *options.RunOptions

	initialized bool
}

var nodeSoftware *NodeSoftware = &NodeSoftware{}

func (s *NodeSoftware) Install() error {
	return nil
}

func (s *NodeSoftware) Exists() bool {
	return true
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
