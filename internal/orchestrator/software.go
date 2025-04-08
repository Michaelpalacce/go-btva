package orchestrator

import (
	"fmt"

	"github.com/Michaelpalacce/go-btva/internal/os/darwin"
	"github.com/Michaelpalacce/go-btva/internal/os/linux"
	"github.com/Michaelpalacce/go-btva/internal/os/software"
)

// installer is a common interface implemented by the installers of all the major systems
type installer interface {
	GetAllSoftware() []software.Software
}

func WithAllSoftware() func(*Orchestrator) error {
	return func(o *Orchestrator) error {
		var installer installer
		switch o.OS.Distro {
		case "linux":
			installer = &linux.Installer{OS: o.OS, Options: o.State.Options}
		case "darwin":
			installer = &darwin.Installer{OS: o.OS, Options: o.State.Options}
		case "windows":
			fallthrough
		default:
			return fmt.Errorf("OS %s is not supported", o.OS.Distro)
		}

		for _, software := range installer.GetAllSoftware() {
			o.SoftwareTasks = append(o.SoftwareTasks, func() error {
				return o.components.softwareComponent.InstallSoftware(software)
			})
		}

		return nil
	}
}
