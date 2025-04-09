package state

import (
	"log/slog"

	"github.com/Michaelpalacce/go-btva/internal/args"
)

// WithCliARgs will load all the arguments from the cli and set them as the options
// @NOTE: This should be used after reading the state from storage
func WithCliArgs() SetStateOption {
	return func(s *State) error {
		if s.Options == nil {
			s.Options = args.Args()
		} else {
			slog.Info("State storage detected and options loaded. Ignoring arguments passed.")
		}

		return nil
	}
}

// @TODO: Expand with other storage if deemed necessary
