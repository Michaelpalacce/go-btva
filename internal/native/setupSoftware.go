package native

import (
	"fmt"
	"log/slog"

	"github.com/Michaelpalacce/go-btva/internal/os/software"
	"github.com/Michaelpalacce/go-btva/internal/state"
)

// installer is a common interface implemented by the installers of all the major systems
type installer interface {
	Java() software.Software
	Mvn() software.Software
	Node() software.Software

	GetAllSoftware() []software.Software
}

// installSoftware is an internal function that can be used to install any software. It will run through a set of commands
// @NOTE: If the version of the software is empty, then we skip installation
func (h *Handler) installSoftware(soft software.Software) error {
	if soft.Exists() || soft.GetVersion() == "" {
		h.state.Set(
			state.WithMsg(soft.GetName(), fmt.Sprintf("Software (%s:%s) already installed, skipping...", soft.GetName(), soft.GetVersion())),
		)
		return nil
	}

	slog.Info("Software is not installed, installing", "name", soft.GetName(), "version", soft.GetVersion())

	if err := soft.Install(); err != nil {
		return err
	}

	h.state.Set(withSoftwareInstalled(soft, nil))

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
