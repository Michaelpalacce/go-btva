package software

import (
	"fmt"

	"github.com/Michaelpalacce/go-btva/pkg/state"
)

const (
	javaStateKey         = "java"
	JAVA_STATE_INSTALLED = "Installed"
)

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
			msg = fmt.Sprintf(JAVA_STATE_INSTALLED)
			step = 1
		}

		s.SetValue(javaStateKey, err == nil, msg, step, err)

		return nil
	}
}

// IsJavaInstalled allows oyu to ask if it isTrue or not
func IsJavaInstalled(isTrue bool) state.GetSuccessStateOption {
	return func(s *state.State) bool {
		value := s.GetValue(javaStateKey)
		if value == nil {
			return !isTrue
		}

		return value.Done == isTrue
	}
}
