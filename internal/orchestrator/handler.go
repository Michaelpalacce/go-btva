package orchestrator

import (
	"github.com/Michaelpalacce/go-btva/internal/args"
	"github.com/Michaelpalacce/go-btva/internal/state"
	"github.com/Michaelpalacce/go-btva/pkg/os"
)

type TaskFunc func() error

// Orchestrator is a struct that orchestrates the setup process based on OS
type Orchestrator struct {
	OS      *os.OS
	State   *state.State
	Options *args.Options

	SoftwareTasks []TaskFunc
	InfraTasks    []TaskFunc
	EnvTasks      []TaskFunc
	FinalTasks    []TaskFunc
}

// NewOrchestrator will return a new Orchestrator that is used to contain and execute tasks
func NewOrchestrator(os *os.OS, state *state.State, options *args.Options) *Orchestrator {
	return &Orchestrator{OS: os, State: state, Options: options}
}

// RunTaskOption accepts a handler and is supposed to modify the state and add tasks to it
type RunTaskOption func(h *Orchestrator) error

// RunTasks executes all the tasks added to the handler in order specified within
// Will resmove all tasks
func (h *Orchestrator) RunTasks(options ...RunTaskOption) error {
	h.SoftwareTasks = make([]TaskFunc, 0)
	h.InfraTasks = make([]TaskFunc, 0)
	h.EnvTasks = make([]TaskFunc, 0)
	h.FinalTasks = make([]TaskFunc, 0)

	for _, option := range options {
		if err := option(h); err != nil {
			return err
		}
	}

	for _, task := range h.SoftwareTasks {
		if err := task(); err != nil {
			return err
		}
	}

	for _, task := range h.InfraTasks {
		if err := task(); err != nil {
			return err
		}
	}

	for _, task := range h.EnvTasks {
		if err := task(); err != nil {
			return err
		}
	}

	for _, task := range h.FinalTasks {
		if err := task(); err != nil {
			return err
		}
	}

	return nil
}
