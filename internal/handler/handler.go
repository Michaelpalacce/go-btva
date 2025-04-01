package handler

import (
	"github.com/Michaelpalacce/go-btva/internal/args"
	"github.com/Michaelpalacce/go-btva/internal/state"
	"github.com/Michaelpalacce/go-btva/pkg/os"
)

type TaskFunc func() error

// Handler is a struct that orchestrates the setup process based on OS
type Handler struct {
	OS      *os.OS
	State   *state.State
	Options *args.Options

	SoftwareTasks []TaskFunc
	InfraTasks    []TaskFunc
	EnvTasks      []TaskFunc
	FinalTasks    []TaskFunc
}

// NewHandler will return a new Handler that will be used to manage and execute os operations
func NewHandler(os *os.OS, state *state.State, options *args.Options) *Handler {
	return &Handler{OS: os, State: state, Options: options}
}

type AddTaskOption func(h *Handler) error

// AddTasks takes a func that works with the Handler and modifies the SoftwareTasks, InfraTasks, EnvTasks and FinalTasks
func (h *Handler) AddTasks(options ...AddTaskOption) error {
	for _, option := range options {
		if err := option(h); err != nil {
			return err
		}
	}

	return nil
}

// RunTasks executes all the tasks added to the handler in order specified within
func (h *Handler) RunTasks() error {
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
