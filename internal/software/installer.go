package software

import (
	"github.com/Michaelpalacce/go-btva/internal/args"
	"github.com/Michaelpalacce/go-btva/internal/software/darwin"
	"github.com/Michaelpalacce/go-btva/internal/software/linux"
	"github.com/Michaelpalacce/go-btva/internal/software/windows"
	"github.com/Michaelpalacce/go-btva/pkg/os"
)

// Installer is a common interface implemented by the installers of all the major systems
type Installer interface {
	InstallNode()
	InstallJava()
	InstallMvn()
	GetOS() *os.OS
}

// NewInstaller will return a pointer to the correct OS installer
func NewInstaller(os *os.OS, options *args.Options) Installer {
	switch os.Distro {
	case "windows":
		return &windows.WindowsInstaller{OS: os, Options: options}
	case "darwin":
		return &darwin.DarwinInstaller{OS: os, Options: options}
	default:
		return &linux.LinuxInstaller{OS: os, Options: options}
	}
}
