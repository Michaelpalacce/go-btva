package run

import (
	"fmt"
	"log/slog"

	"github.com/Michaelpalacce/go-btva/internal/run/orchestrator"
	"github.com/Michaelpalacce/go-btva/internal/run/state"
	osl "github.com/Michaelpalacce/go-btva/pkg/os"
)

type RunCommand struct{}

func (c *RunCommand) Name() string {
	return "run"
}

func (c *RunCommand) Run() error {
	var (
		o     *orchestrator.Orchestrator
		err   error
		osPtr *osl.OS
		s     *state.State
	)

	// Init Block. Used for fetching and creating needed structs

	if s, err = state.NewState(state.WithDefaultJsonStorage(), c.WithRunCliArgs()); err != nil {
		slog.Error("Error while loading state.", "err", err)
		return fmt.Errorf("error while loading state. Err %w", err)
	}

	osPtr = osl.GetOS()

	o = orchestrator.NewOrchestrator(osPtr, s)

	err = o.Tasks(orchestrator.WithOptions())
	if err != nil {
		return fmt.Errorf("error while adding tasks. Err %w", err)
	}

	err = o.Run()
	if err != nil {
		return fmt.Errorf("error while running tasks. Err %w", err)
	}

	return nil
}

// WithCliARgs will load all the arguments from the cli and set them as the options
// @NOTE: This should be used after reading the state from storage
func (c *RunCommand) WithRunCliArgs() state.SetStateOption {
	return func(s *state.State) error {
		if s.Options == nil {
			s.Options = c.Args()
		} else {
			slog.Info("State storage detected and options loaded. Ignoring arguments passed.")
		}

		return nil
	}
}
