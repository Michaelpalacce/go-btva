package software

import (
	"fmt"
	"log/slog"

	"github.com/Michaelpalacce/go-btva/internal/args"
	"github.com/Michaelpalacce/go-btva/internal/handler"
	"github.com/Michaelpalacce/go-btva/internal/os/darwin"
	"github.com/Michaelpalacce/go-btva/internal/os/linux"
	"github.com/Michaelpalacce/go-btva/internal/os/software"
	"github.com/Michaelpalacce/go-btva/internal/state"
	"github.com/Michaelpalacce/go-btva/pkg/os"
)

// installer is a common interface implemented by the installers of all the major systems
type installer interface {
	GetAllSoftware() []software.Software
}

type SoftwareComponent struct {
	os      *os.OS
	state   *state.State
	options *args.Options
}

func NewSoftware(os *os.OS, state *state.State) *SoftwareComponent {
	return &SoftwareComponent{os: os, state: state, options: state.Options}
}

func WithAllSoftware() func(*handler.Handler) error {
	return func(h *handler.Handler) error {
		softwareComponent := NewSoftware(h.OS, h.State)

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

// installSoftware is an internal function that can be used to install any software. It will run through a set of commands
// @NOTE: If the version of the software is empty, then we skip installation
func (s *SoftwareComponent) InstallSoftware(soft software.Software) error {
	if soft.Exists() || soft.GetVersion() == "" {
		s.state.Set(
			state.WithWarn(soft.GetName(), fmt.Sprintf("Software (%s:%s) already installed, skipping...", soft.GetName(), soft.GetVersion())),
		)
		return nil
	}

	slog.Info("Software is not installed, installing", "name", soft.GetName(), "version", soft.GetVersion())

	if err := soft.Install(); err != nil {
		s.state.Set(withSoftwareInstalled(soft, err))
		return err
	}

	s.state.Set(withSoftwareInstalled(soft, nil))

	return nil
}

// withSoftwareInstalled will set the state of the given software
func withSoftwareInstalled(soft software.Software, err error) state.SetStateOption {
	return func(s *state.State) error {
		var msg string
		if err != nil {
			msg = fmt.Sprintf("Error installing %s:%s. Error was %v", soft.GetName(), soft.GetVersion(), err)
		} else {
			msg = fmt.Sprintf("Software %s:%s was installed", soft.GetName(), soft.GetVersion())
		}

		s.Set(
			state.WithMsg(soft.GetName(), msg),
			state.WithErr(soft.GetName(), err),
		)

		return nil
	}
}
