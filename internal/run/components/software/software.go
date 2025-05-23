package software

import (
	"fmt"
	"log/slog"

	"github.com/Michaelpalacce/go-btva/internal/options"
	"github.com/Michaelpalacce/go-btva/internal/run/os/software"
	"github.com/Michaelpalacce/go-btva/internal/run/state"
	"github.com/Michaelpalacce/go-btva/pkg/os"
	"github.com/Michaelpalacce/go-btva/pkg/prompt"
)

type SoftwareComponent struct {
	os      *os.OS
	state   *state.State
	options *options.RunOptions
}

func NewSoftwareComponent(os *os.OS, state *state.State) *SoftwareComponent {
	return &SoftwareComponent{os: os, state: state, options: state.Options}
}

// installSoftware is an internal function that can be used to install any software. It will run through a set of commands
// @NOTE: If the version of the software is empty, then we skip installation
func (s *SoftwareComponent) InstallSoftware(soft software.Software) error {
	if soft.Exists() || soft.GetVersion() == "" {
		version := soft.GetVersion()
		if version != "" {
			version = ":" + version
		}
		s.state.Set(
			state.WithWarn(soft.GetName(), fmt.Sprintf("Software (%s%s) already installed, skipping...", soft.GetName(), version)),
		)
		return nil
	}

	if s.options.Prompt {
		if a, err := prompt.AskYesNoAbort(fmt.Sprintf("%s:%s will be installed.", soft.GetName(), soft.GetVersion())); err != nil {
			return err
		} else if !a {
			return nil
		}
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
