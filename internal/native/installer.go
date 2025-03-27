package native

import (
	"fmt"
	"log/slog"

	"github.com/Michaelpalacce/go-btva/internal/software"
	"github.com/Michaelpalacce/go-btva/pkg/state"
)

// Installer is a common interface implemented by the installers of all the major systems
type Installer interface {
	Java() software.Software
	Mvn() software.Software
	Node() software.Software
}

// installSoftware is an internal function that can be used to install any software. It will run through a set of commands
func (h *Handler) installSoftware(soft software.Software) error {
	if h.state.GetDone(softwareDone(soft)) {
		slog.Info("Software already installed, skipping...", "name", soft.GetName(), "version", soft.GetVersion())
		return nil
	}

	if soft.Exists() {
		h.state.Set(withSoftwareInstalled(soft, nil))
		return nil
	}

	slog.Info("Software is not installed, installing", "name", soft.GetName(), "version", soft.GetVersion())

	if err := soft.Install(); err != nil {
		h.state.Set(withSoftwareInstalled(soft, err))
		return err
	}

	h.state.Set(withSoftwareInstalled(soft, nil))

	return nil
}

// withSoftwareInstalled will set the state of the given software
func withSoftwareInstalled(soft software.Software, err error) state.SetStateOption {
	return func(s *state.State) error {
		var (
			msg  string
			step int
		)
		if err != nil {
			msg = fmt.Sprintf("Error installing %s:%s. Error was %v", soft.GetName(), soft.GetVersion(), err)
			step = 0
		} else {
			msg = fmt.Sprintf("Software %s:%s was installed", soft.GetName(), soft.GetVersion())
			step = 1
		}

		s.Set(
			state.WithDone(soft.GetName(), err == nil),
			state.WithMsg(soft.GetName(), msg),
			state.WithErr(soft.GetName(), err),
			state.WithStep(soft.GetName(), step),
		)

		return nil
	}
}

// softwareDone checks if the current software is Done
func softwareDone(soft software.Software) state.GetSuccessStateOption {
	return func(s *state.State) bool {
		value := s.GetValue(soft.GetName())
		if value == nil {
			return false
		}

		return value.Done
	}
}
