package software

import (
	"fmt"

	"github.com/Michaelpalacce/go-btva/pkg/state"
)

const (
	SOFTWARE_STATE_INSTALLED = "Installed"
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
			msg = fmt.Sprintf(SOFTWARE_STATE_INSTALLED)
			step = 1
		}

		s.SetValue(software.GetName(), err == nil, msg, step, err)

		return nil
	}
}

// IsJavaInstalled allows oyu to ask if it isTrue or not
func IsJavaInstalled(isTrue bool) state.GetSuccessStateOption {
	return func(s *state.State) bool {
		value := s.GetValue(JavaSoftwareKey)
		if value == nil {
			return !isTrue
		}

		return value.Done == isTrue
	}
}

// IsMvnInstalled allows you to ask if it isTrue or not
func IsMvnInstalled(isTrue bool) state.GetSuccessStateOption {
	return func(s *state.State) bool {
		value := s.GetValue(MvnSoftwareKey)
		if value == nil {
			return !isTrue
		}

		return value.Done == isTrue
	}
}
