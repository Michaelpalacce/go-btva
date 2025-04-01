package orchestrator

import (
	"fmt"

	software_component "github.com/Michaelpalacce/go-btva/internal/components/software"
	"github.com/Michaelpalacce/go-btva/internal/os/darwin"
	"github.com/Michaelpalacce/go-btva/internal/os/linux"
	"github.com/Michaelpalacce/go-btva/internal/os/software"
)

// installer is a common interface implemented by the installers of all the major systems
type installer interface {
	GetAllSoftware() []software.Software
}

func WithAllSoftware() func(*Orchestrator) error {
	return func(h *Orchestrator) error {
		softwareComponent := software_component.NewSoftware(h.OS, h.State)

		var installer installer
		switch h.OS.Distro {
		case "linux":
			installer = &linux.Installer{OS: h.OS, Options: h.Options}
		case "darwin":
			installer = &darwin.Installer{OS: h.OS, Options: h.Options}
		case "windows":
			fallthrough
		default:
			return fmt.Errorf("OS %s is not supported", h.OS.Distro)
		}

		for _, software := range installer.GetAllSoftware() {
			h.SoftwareTasks = append(h.SoftwareTasks, func() error {
				return softwareComponent.InstallSoftware(software)
			})
		}

		return nil
	}
}
