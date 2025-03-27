package software

import (
	"fmt"

	"github.com/Michaelpalacce/go-btva/pkg/state"
)

// SoftwareInstalled will set the state of the given software
func SoftwareInstalled(software Software, err error) state.SetStateOption {
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

		s.SetValue(software.GetName(), err == nil, msg, step, err)

		return nil
	}
}

// IsSoftwareInstalled checks if the current software is installed
func IsSoftwareInstalled(software Software) state.GetSuccessStateOption {
	return func(s *state.State) bool {
		value := s.GetValue(software.GetName())
		if value == nil {
			return false
		}

		return value.Done
	}
}

// IsSoftwareNotInstalled checks if the current software is not installed
func IsSoftwareNotInstalled(software Software) state.GetSuccessStateOption {
	return func(s *state.State) bool {
		value := s.GetValue(software.GetName())
		if value == nil {
			return true
		}

		return !value.Done
	}
}
