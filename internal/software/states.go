package software

import (
	"fmt"

	"github.com/Michaelpalacce/go-btva/pkg/state"
)

// WithSoftwareInstalled will set the state of the given software
func WithSoftwareInstalled(software Software, err error) state.SetStateOption {
	return func(s *state.State) error {
		var (
			msg  string
			step int
		)
		if err != nil {
			msg = fmt.Sprintf("Error installing %s:%s. Error was %v", software.GetName(), software.GetVersion(), err)
			step = 0
		} else {
			msg = fmt.Sprintf("Software %s:%s was installed", software.GetName(), software.GetVersion())
			step = 1
		}

		s.Set(
			state.WithDone(software.GetName(), err == nil),
			state.WithMsg(software.GetName(), msg),
			state.WithErr(software.GetName(), err),
			state.WithStep(software.GetName(), step),
		)

		return nil
	}
}

// SoftwareDone checks if the current software is Done
func SoftwareDone(software Software) state.GetSuccessStateOption {
	return func(s *state.State) bool {
		value := s.GetValue(software.GetName())
		if value == nil {
			return false
		}

		return value.Done
	}
}
