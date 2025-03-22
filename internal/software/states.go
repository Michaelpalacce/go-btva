package software

import (
	"fmt"

	"github.com/Michaelpalacce/go-btva/pkg/state"
)

const javaInstallationState = "javaSoftware"

// JAVA
func JavaInstalled(err error) state.SetStateOption {
	return func(s *state.State) error {
		var (
			msg  string
			step int
		)
		if err != nil {
			msg = fmt.Sprintf("Error installing Java: %v", err)
			step = 0
		} else {
			msg = fmt.Sprintf("Java installed")
			step = 1
		}

		s.SetValue(javaInstallationState, err != nil, msg, step, err)

		return nil
	}
}

func IsJavaInstalled(isTrue bool) state.GetSuccessStateOption {
	return func(s *state.State) bool {
		value := s.GetValue(javaInstallationState)
		if value == nil {
			return isTrue && false
		}

		return value.Done && isTrue
	}
}
