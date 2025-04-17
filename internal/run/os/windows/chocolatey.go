package windows

import (
	"github.com/Michaelpalacce/go-btva/internal/options"
	"github.com/Michaelpalacce/go-btva/internal/run/os/software"
	"github.com/Michaelpalacce/go-btva/pkg/command/windows"
	"github.com/Michaelpalacce/go-btva/pkg/os"
)

// ChocolateySoftware is responsible for installing and checking if java is installed
type ChocolateySoftware struct {
	os      *os.OS
	options *options.RunOptions

	initialized bool
}

var chocolateySoftware *ChocolateySoftware = &ChocolateySoftware{}

func (s *ChocolateySoftware) Install() error {
	if err := windows.RunElevatedCommand(`Set-ExecutionPolicy Bypass -Scope Process -Force; [System.Net.ServicePointManager]::SecurityProtocol = [System.Net.ServicePointManager]::SecurityProtocol -bor 3072; iex ((New-Object System.Net.WebClient).DownloadString('https://community.chocolatey.org/install.ps1'))`); err != nil {
		return err
	}

	return nil
}

func (s *ChocolateySoftware) Exists() bool {
	return false
}

func (s *ChocolateySoftware) GetName() string    { return software.ChocolateySoftwareKey }
func (s *ChocolateySoftware) GetVersion() string { return s.options.Software.ChocolateyVersion }

// Chocolatey will return the ChocolateySoftware object that can be used to install and check if java exists
// Only a single instance of the ChocolateySoftware will be returned
func (i *Installer) Chocolatey() software.Software {
	if !chocolateySoftware.initialized {
		chocolateySoftware.os = i.OS
		chocolateySoftware.options = i.Options
		chocolateySoftware.initialized = true
	}

	return chocolateySoftware
}
